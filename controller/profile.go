package controller

import (
	"backend/utils"
	"context"
	"encoding/json"
	"net/http"
)

// GetUserProfileHandler retrieves the user profile from Firebase Realtime Database
func GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the UID from the request context (set by the middleware)
	uid := r.Context().Value("uid").(string)
	if uid == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Fetch user details from Firebase Database
	var userProfile map[string]interface{}
	err := utils.FirebaseDB.NewRef("users/"+uid).Get(context.Background(), &userProfile)
	if err != nil {
		http.Error(w, "Failed to retrieve user profile", http.StatusInternalServerError)
		return
	}

	// Return the user's profile as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(userProfile); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
