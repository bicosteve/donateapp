package models

import (
	"context"
	"golang.org/x/crypto/bcrypt"
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

func (u *User) findByEmail(user User) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	q := "SELECT email FROM users WHERE email = ?"
	row := db.QueryRowContext(ctx, q, user.Email)
	err := row.Scan(&user.Email)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Password hashing
func (u *User) hashPassword(user User) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (u *User) RegisterUser(user User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	isUser, err := u.findByEmail(user)
	if isUser == true {
		return nil, err
	}

	hashedPassword, err := u.hashPassword(user)
	if err != nil {
		return nil, err
	}

	q := "INSERT INTO users (email, phone_number,password,created_at,updated_at) VALUES (?,?,?,?,?)"

	// EXECUTE QUERY
	_, err = db.ExecContext(
		ctx, q, user.Email, user.PhoneNumber, hashedPassword, time.Now(), time.Now())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	
	return &user, nil

}

func passwordCompare() {}
