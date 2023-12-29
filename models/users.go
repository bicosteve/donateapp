package models

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"path/filepath"
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

type Claims struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	jwt.RegisteredClaims
}

func (u *User) RegisterUser(user User) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := u.hashPassword(user)
	if err != nil {
		return nil, err
	}

	q := "INSERT INTO users (email, phone_number,password,created_at,updated_at) VALUES (?,?,?,?,?)"

	_, err = db.ExecContext(
		ctx, q, user.Email, user.PhoneNumber, hashedPassword, time.Now(), time.Now())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &user, nil

}

func (u *User) GenerateAuthToken(user User) (string, error) {
	//fileName := ".env"
	path, err := filepath.Abs(".env")
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load(filepath.Join(path))
	if err != nil {
		return "", errors.New("cannot load .env for jwt")
	}

	jwtKey := os.Getenv("JWTSECRET")
	claims := &Claims{
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (u *User) FindByEmail(user User) (bool, error) {
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

func (u *User) hashPassword(user User) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("something went wrong with password hashing")
	}
	return string(bytes), nil
}

func (u *User) PasswordCompare(user User) bool {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var dbUser User
	q := "SELECT * FROM users WHERE email = ?"
	row := db.QueryRowContext(ctx, q, user.Email)
	err := row.Scan(
		&dbUser.ID,
		&dbUser.Email,
		&dbUser.PhoneNumber,
		&dbUser.Password,
		&dbUser.CreateAt,
		&dbUser.UpdatedAt,
	)

	if err != nil {
		log.Fatal(err)
		return false
	}

	userPassword := []byte(user.Password)
	dbPassword := []byte(dbUser.Password)

	err = bcrypt.CompareHashAndPassword(dbPassword, userPassword)
	if err != nil {
		return false
	}
	return true
}
