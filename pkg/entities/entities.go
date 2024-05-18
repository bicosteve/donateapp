package entities

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type User struct {
	ID              string    `json:"id"`
	Email           string    `json:"email"`
	PhoneNumber     string    `json:"phone_number"`
	Password        string    `json:"password"`
	ConfirmPassword string    `json:"confirm_password"`
	CreateAt        time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type UserRequestBody struct {
	Email           string `json:"email"`
	PhoneNumber     string `json:"phone_number"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type Claims struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	jwt.RegisteredClaims
}

type Donation struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Photo     string    `json:"photo"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    string    `json:"userid"`
}

type DonationBody struct {
	Name     string `json:"name"`
	Photo    string `json:"photo"`
	Location string `json:"location"`
}
