package controllers

import (
	"donateapp/helpers"
	"donateapp/models"
	"encoding/json"
	"net/http"
	"time"
)

// POST User -> api/v1/user/register
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	isValidNumber := helpers.CheckPhoneNumber(user)
	if isValidNumber == false {
		msg := "Phone number must be 10 numbers"
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": msg})
		helpers.MessageLogs.ErrorLog.Println(msg)
		return
	}

	isValidEmail := helpers.IsValidEmail(user.Email)
	if isValidEmail == false {
		msg := "Provide valid email address"
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": msg})
		helpers.MessageLogs.ErrorLog.Println(msg)
		return
	}

	isValidPassword := helpers.ValidatePassword(user)
	if isValidPassword == false {
		msg := "Password cannot be empty & must match confirm password"
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": msg})
		helpers.MessageLogs.ErrorLog.Println(msg)
		return
	}

	isUser, err := user.FindByEmail(user)
	if isUser == true {
		msg := "User already exists"
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": msg})
		//helpers.ErrorJSON(w, err, http.StatusBadRequest)
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	createdUser, err := user.RegisterUser(user)
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
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	isValidEmail := helpers.IsValidEmail(user.Email)
	if isValidEmail == false {
		msg := "Provide valid email address"
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": msg})
		helpers.MessageLogs.ErrorLog.Println(msg)
		return
	}

	isUser, err := user.FindByEmail(user)
	if isUser == false {
		msg := "User does not exist"
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": msg})
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	isValidPassword := user.PasswordCompare(user)
	if isValidPassword == false {
		msg := "Password and confirm password do not match"
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": msg})
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	token, err := user.GenerateAuthToken(user)

	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, helpers.Envelope{"msg": err})
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	// Setting the cookie on headers
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(time.Hour * 1),
	})

}

//func GetProfile(w http.ResponseWriter, r *http.Request) {
//	cookie, err := r.Cookie("token")
//	if err != nil {
//		if err == http.ErrNoCookie {
//			helpers.WriteJSON(w, http.StatusUnauthorized, helpers.Envelope{"msg": "Forbidden"})
//			return
//		}
//
//		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": "Bad Request"})
//		return
//	}
//
//	tokenStr := cookie.Value
//
//	claims := &models.Claims{}
//
//	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
//		func(t *jwt.Token) (interface{}, error) {
//			return jwtKey, nil
//		})
//}
