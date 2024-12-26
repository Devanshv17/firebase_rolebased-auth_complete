package controller

import (
	"backend/model"
	"backend/utils"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// LoginHandler generates token and sends it to the client
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(body))

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

	// Retrieve hashed password from Firebase Database
	var hashedPassword string
	err = utils.FirebaseDB.NewRef("users/"+u.UID+"/hashed_password").Get(context.Background(), &hashedPassword)
	if err != nil || hashedPassword == "" {
		http.Error(w, "Failed to retrieve user password", http.StatusInternalServerError)
		return
	}

	// Trim the input password to avoid issues with whitespace
	inputPassword := strings.TrimSpace(user.Password)

	// Compare stored hashed password with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Retrieve user's role from Firebase Database
	var role string
	err = utils.FirebaseDB.NewRef("users/"+u.UID+"/role").Get(context.Background(), &role)
	if err != nil || role == "" {
		http.Error(w, "Failed to retrieve user role", http.StatusInternalServerError)
		return
	}

	// Generate JWT token for the user with UID and role
	token, err := utils.GenerateJWT(u.UID, role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Prepare response payload with UID, role, and the generated token
	response := map[string]interface{}{
		"user_id":   u.UID,
		"role":      role,
		"jwt_token": token, // Send JWT token back to the client
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		log.Printf("Failed to encode response: %v\n", err)
	}
}
