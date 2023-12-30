package controllers

import (
	"donateapp/helpers"
	"donateapp/models"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// POST User -> api/v1/user/register
func CreateUser(w http.ResponseWriter, r *http.Request) {
	userReqBody := new(models.UserRequestBody)

	err := json.NewDecoder(r.Body).Decode(&userReqBody)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	isValidNumber := helpers.CheckPhoneNumber(userReqBody.PhoneNumber)
	if isValidNumber == false {
		msg := "Phone number must be 10 numbers"
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": msg})
		helpers.MessageLogs.ErrorLog.Println(msg)
		return
	}

	isValidEmail := helpers.IsValidEmail(userReqBody.Email)
	if isValidEmail == false {
		msg := "Provide valid email address"
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": msg})
		helpers.MessageLogs.ErrorLog.Println(msg)
		return
	}

	isValidPassword := helpers.ValidatePassword(userReqBody.Password, userReqBody.ConfirmPassword)
	if isValidPassword == false {
		msg := "Password cannot be empty & must match confirm password"
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": msg})
		helpers.MessageLogs.ErrorLog.Println(msg)
		return
	}

	isUser, err := userReqBody.FindByEmail(userReqBody.Email)
	if isUser == true {
		msg := "User already exists"
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": msg})
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	createdUser, err := userReqBody.RegisterUser(*userReqBody)
	if err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": err})
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	helpers.WriteJSON(w, http.StatusCreated, createdUser)
}

// Login POST -> /api/users/login
func LoginUser(w http.ResponseWriter, r *http.Request) {
	userReqBody := new(models.UserRequestBody)
	//user := new(models.User)
	err := json.NewDecoder(r.Body).Decode(&userReqBody)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	isValidEmail := helpers.IsValidEmail(userReqBody.Email)
	if isValidEmail == false {
		msg := "Invalid email address"
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": msg})
		helpers.MessageLogs.ErrorLog.Println(msg)
		return
	}

	userExists, err := userReqBody.FindByEmail(userReqBody.Email)
	if userExists == false {
		msg := "User does not exist"
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": msg})
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	isValidPassword := userReqBody.PasswordCompare(*userReqBody)
	if isValidPassword == false {
		msg := "Password and confirm password do not match"
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": msg})
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	tokenString, err := userReqBody.GenerateAuthToken(*userReqBody)
	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, helpers.Envelope{"msg": err})
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	// Setting the cookie on headers
	cookie := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		SameSite: 4,
		Secure:   true,
		Expires:  time.Now().Add(time.Second * 20),
	}
	http.SetCookie(w, cookie)
	r.AddCookie(cookie)
}

// Get user profile -> GET -> /api/users/profile
func GetProfile(w http.ResponseWriter, r *http.Request) {
	path, err := filepath.Abs(".env")
	if err != nil {
		log.Fatal(err)
	}
	err = godotenv.Load(filepath.Join(path))
	if err != nil {
		log.Fatal(err)
	}
	jwtKey := os.Getenv("JWTSECRET")

	cookie, err := r.Cookie("token")
	if err != nil {
		fmt.Println(err)
		if err == http.ErrNoCookie {
			helpers.WriteJSON(w, http.StatusUnauthorized, helpers.Envelope{"msg": "Forbidden. No Cookie"})
			return
		}
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": "Bad Request"})
		return
	}

	tokenStr := cookie.Value
	claims := &models.Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			helpers.WriteJSON(w, http.StatusUnauthorized, helpers.Envelope{"msg": "Invalid token"})
			return
		}

		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": "Something went wrong"})
		return

	}

	if !tkn.Valid {
		helpers.WriteJSON(w, http.StatusUnauthorized, helpers.Envelope{"msg": "Unauthorized"})
		return
	}

	fmt.Println([]byte(claims.Email))
}
