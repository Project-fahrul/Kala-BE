package entity

import "time"

type Users struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Name         string    `json:"name"`
	PhoneNumber  string    `json:"phone_number"`
	Role         string    `json:"role"`
	Token        string    `json:"token"`
	TokenExpired time.Time `json:"token_expired"`
	LoginDelay   time.Time `json:"login_delay"`
	Verified     bool      `json:"verified"`
}
