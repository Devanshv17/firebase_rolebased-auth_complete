package main

import (
	"backend/controller"
	"backend/middleware"
	"backend/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize Firebase Auth and Database clients
	utils.InitFirebase()

	r := mux.NewRouter()

	// Apply CORS middleware globally
	r.Use(middleware.CORS)

	// Register routes that do not require authentication
	r.HandleFunc("/register", controller.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", controller.LoginHandler).Methods("POST")
	r.HandleFunc("/forget-password", controller.ForgotPasswordHandler).Methods("POST")
	r.HandleFunc("/resend-verification", controller.ResendVerificationHandler).Methods("POST")

	// Apply AuthMiddleware to routes that require authentication
	authenticatedRoutes := r.PathPrefix("/user").Subrouter()
	authenticatedRoutes.Use(middleware.AuthMiddleware)
	authenticatedRoutes.HandleFunc("/profile", controller.GetUserProfileHandler).Methods("GET")
	authenticatedRoutes.HandleFunc("/enter_data", controller.EnterDataHandler).Methods("POST")

	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
