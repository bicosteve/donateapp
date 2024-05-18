package controllers

import (
	"donateapp/pkg/helpers"
	"donateapp/pkg/models"
	"net/http"
	"strconv"
	"time"
)

// POST User -> api/v1/user/register
func CreateUser(w http.ResponseWriter, r *http.Request) {
	userReqBody := new(models.UserRequestBody)

	err := helpers.ReadJSON(w, r, userReqBody)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
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

// Login POST -> api/users/login
func LoginUser(w http.ResponseWriter, r *http.Request) {
	userReqBody := new(models.UserRequestBody)

	err := helpers.ReadJSON(w, r, userReqBody)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
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
		Expires:  time.Now().Add(time.Hour * 24),
	}
	http.SetCookie(w, cookie)
	r.AddCookie(cookie)
}

// Get user profile -> GET -> /api/users/profile
func GetProfile(w http.ResponseWriter, r *http.Request) {
	jwtKey, err := helpers.LoadJWTKEY()
	if err != nil {
		helpers.WriteJSON(w, http.StatusUnauthorized, helpers.Envelope{"msg": err})
		return
	}

	tokenString, err := helpers.GenerateTokenString(r)
	if err != nil {
		helpers.WriteJSON(w, http.StatusUnauthorized, helpers.Envelope{"msg": err})
		return
	}

	claims, err := helpers.ValidClaim(&models.Claims{}, tokenString, jwtKey)
	if err != nil {
		helpers.WriteJSON(w, http.StatusUnauthorized, helpers.Envelope{"msg": err})
		return
	}

	var user models.User
	userID, _ := strconv.Atoi(claims.ID)
	userProfile, err := user.GetProfile(userID)

	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, helpers.Envelope{"msg": "Internal Server Error"})
		return
	}

	helpers.WriteJSON(w, http.StatusOK, userProfile)
}
