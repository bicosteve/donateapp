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

	token, err := userReqBody.GenerateAuthToken(*userReqBody)

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
