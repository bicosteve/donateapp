package models

import (
	"context"
	"donateapp/helpers"
	"log"
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

func (u *User) RegisterUser(user User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)

	defer cancel()

	hashPassword, err := helpers.HashPassword(user.Password)

	if err != nil {
		return nil, err
	}

	registerQuery := `
						INSERT INTO users (email, phone_number,password,created_at,updated_at)
						VALUES ($1,$2,$3,$4,$5) 
						`

	// EXECUTE QUERY
	_, err = db.ExecContext(
		ctx, registerQuery, user.Email, user.PhoneNumber, hashPassword, time.Now(), time.Now())

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &user, nil

}
