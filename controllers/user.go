package controllers

import (
	"donateapp/models"
	"fmt"
	"net/http"
)

var user models.User

// POST User -> api/v1/user/register

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	fmt.Println(user)

}
