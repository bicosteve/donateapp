package controllers

import (
	"database/sql"
	"donateapp/helpers"
	"donateapp/models"
	"encoding/json"
	"net/http"
)

var db *sql.DB

var user models.User

// POST User -> api/v1/user/register

func CreateUser(w http.ResponseWriter, r *http.Request) {
	//ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	var user models.User

	defer db.Close()

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println(err)
	}

	isValidUser := helpers.CheckValidUser(user)
	if isValidUser == false {
		msg := "All fields are required"
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
		msg := "Password and confirm password do not match"
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": msg})
		helpers.MessageLogs.ErrorLog.Println(msg)
		return
	}
	//
	//found, err := helpers.UserExists(db, user)
	//
	//if found == true {
	//	msg := "User already exists"
	//	helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": msg})
	//	helpers.MessageLogs.ErrorLog.Println(msg)
	//	return
	//}

	createdUser, err := user.RegisterUser(user)
	if err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": err})
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusCreated, createdUser)
}
