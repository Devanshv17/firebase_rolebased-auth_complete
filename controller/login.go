package controller

import (
	"backend/model"
	"backend/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// LoginHandler generates token and sends it to the client
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Password == "" {
		http.Error(w, "Email and Password are required", http.StatusBadRequest)
		return
	}

	// Authenticate user by email using Firebase Auth
	u, err := utils.FirebaseAuth.GetUserByEmail(context.Background(), user.Email)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if !u.EmailVerified {
		http.Error(w, "Email not verified", http.StatusUnauthorized)
		return
	}

	// Retrieve user details (hashed_password and role) in a single Firebase call
	type UserDetails struct {
		HashedPassword string `json:"hashed_password"`
		Role           string `json:"role"`
	}

	var userDetails UserDetails
	err = utils.FirebaseDB.NewRef("users/"+u.UID).Get(context.Background(), &userDetails)
	if err != nil || userDetails.HashedPassword == "" || userDetails.Role == "" {
		http.Error(w, "Failed to retrieve user details", http.StatusInternalServerError)
		return
	}

	// Compare stored hashed password with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(userDetails.HashedPassword), []byte(user.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token for the user with UID and role
	token, err := utils.GenerateJWT(u.UID, userDetails.Role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Prepare response payload
	response := map[string]interface{}{
		"role":      userDetails.Role,
		"jwt_token": token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		log.Printf("Failed to encode response: %v\n", err)
	}
}
