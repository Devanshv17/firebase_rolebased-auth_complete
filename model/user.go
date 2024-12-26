package model

type User struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Role        string `json:"role"`         // 'user' or 'organizer'
	PhoneNumber string `json:"phone_number"` // Unique phone number
	Name        string `json:"name"`
	Gender      string `json:"gender"` // 'male', 'female', or 'others'
	City        string `json:"city"`
}
