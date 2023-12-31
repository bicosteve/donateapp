package controllers

import (
	"donateapp/helpers"
	"donateapp/models"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

var donation models.Donation

// POST -> /api/v1/donations/donation

func CreateDonation(w http.ResponseWriter, r *http.Request) {
	var donationData models.Donation
	err := json.NewDecoder(r.Body).Decode(&donationData)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println("Error", err)
		return
	}

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
		if err == http.ErrNoCookie {
			helpers.WriteJSON(w, http.StatusUnauthorized, helpers.Envelope{"msg": "Forbidden. No Cookie"})
			return
		}

		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": "Bad Request"})
		return
	}

	tokenString := cookie.Value
	claims := &models.Claims{}
	tkn, err := jwt.ParseWithClaims(tokenString, claims,
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

	userID, _ := strconv.Atoi(claims.ID)
	donation, err := donation.AddDonations(donationData, userID)

	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, helpers.Envelope{"msg": "Cannot add donations"})
		return
	}

	helpers.WriteJSON(w, http.StatusCreated, donation)

}
