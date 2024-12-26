package controller

import (
	"backend/model"
	"backend/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"firebase.google.com/go/auth"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Ensure that the role is not empty
	if user.Role == "" {
		http.Error(w, "Role is required", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		log.Printf("Failed to hash password: %v\n", err)
		return
	}

	// Create the Firebase user
	params := (&auth.UserToCreate{}).
		Email(user.Email).
		Password(string(hashedPassword)).
		DisplayName(user.Role) // Ensure role is passed here as the DisplayName

	newUser, err := utils.FirebaseAuth.CreateUser(context.Background(), params)
	if err != nil {
		http.Error(w, "Email already Exists", http.StatusInternalServerError)
		log.Printf("Failed to create user: %v\n", err)
		return
	}

	// Send verification email
	err = utils.SendVerificationEmail(newUser)
	if err != nil {
		http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
		log.Printf("Failed to send verification email: %v\n", err)
		return
	}

	// Assign role to the user in Firebase Database
	err = utils.FirebaseDB.NewRef("users/"+newUser.UID+"/role").Set(context.Background(), user.Role)
	if err != nil {
		http.Error(w, "Failed to assign role to user", http.StatusInternalServerError)
		log.Printf("Failed to assign role to user: %v\n", err)
		return
	}

	// Save hashed password in Firebase Database
	err = utils.FirebaseDB.NewRef("users/"+newUser.UID+"/hashed_password").Set(context.Background(), string(hashedPassword))
	if err != nil {
		http.Error(w, "Failed to save user password", http.StatusInternalServerError)
		log.Printf("Failed to save user password: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}
