package utils

import (
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Secret key used to sign the JWTs
var jwtKey = []byte("your_secret_key_here")

// Claims struct for JWT payload
type Claims struct {
	UID  string `json:"uid"`
	Role string `json:"role"`
	jwt.StandardClaims
}

// GenerateJWT generates a JWT token for the user
func GenerateJWT(uid, role string) (string, error) {
	// Set expiration time for the token
	expirationTime := time.Now().Add(24 * time.Hour) // 24 hours expiration

	// Create JWT claims
	claims := &Claims{
		UID:  uid,
		Role: role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "Zintrix", // Change this to your app's name
		},
	}

	// Create the token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token using the secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return "", err
	}

	return tokenString, nil
}

// VerifyToken verifies the JWT token and returns the claims
func VerifyToken(tokenString string) (*Claims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return nil, err
	}

	// Check if the token is valid
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		log.Println("Invalid token")
		return nil, err
	}

	return claims, nil
}
