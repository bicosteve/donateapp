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

func (u *UserRequestBody) RegisterUser(user UserRequestBody) (*UserRequestBody, error) {
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

func (u *UserRequestBody) FindByEmail(email string) (bool, error) {
	var user = new(User)
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	q := "SELECT id, email  FROM users WHERE email = ?"
	row := db.QueryRowContext(ctx, q, email)
	err := row.Scan(&user.ID, &user.Email)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *UserRequestBody) hashPassword(user UserRequestBody) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("something went wrong with password hashing")
	}
	return string(bytes), nil
}

func (u *UserRequestBody) PasswordCompare(user UserRequestBody) bool {
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

func (u *UserRequestBody) GenerateAuthToken(user UserRequestBody) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	path, err := filepath.Abs(".env")
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load(filepath.Join(path))
	if err != nil {
		return "", errors.New("cannot load .env for jwt")
	}

	// Getting the user from db
	var dbUser User
	query := `SELECT * FROM users WHERE email = ?`
	row := db.QueryRowContext(ctx, query, user.Email)
	err = row.Scan(
		&dbUser.ID,
		&dbUser.Email,
		&dbUser.PhoneNumber,
		&dbUser.Password,
		&dbUser.CreateAt,
		&dbUser.UpdatedAt,
	)

	if err != nil {
		log.Fatal(err)
	}

	jwtKey := os.Getenv("JWTSECRET")
	claims := &Claims{
		ID:          dbUser.ID,
		Email:       dbUser.Email,
		PhoneNumber: dbUser.PhoneNumber,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (u *User) GetProfile(userId int) (*User, error) {
	//var user = new(User)
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	q := `SELECT id, email,phone_number  FROM users WHERE id = ?`
	row := db.QueryRowContext(ctx, q, userId)
	err := row.Scan(&user.ID, &user.Email, &user.PhoneNumber)

	if err != nil {
		log.Fatal(err)
	}
	return &user, nil

}
