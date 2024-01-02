package controllers

import (
	"donateapp/helpers"
	"donateapp/models"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

var donation models.Donation
var donationData models.Donation

// POST -> /api/v1/donations/donation
func CreateDonation(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&donationData)
	if err != nil {
		helpers.MessageLogs.ErrorLog.Println("Error", err)
		return
	}

	isValid := helpers.ValidateDonationPayload(donationData)
	if !isValid {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": "All fields are required"})
		return
	}

	jwtKey, err := helpers.LoadJWTKEY() // Load JWT Key
	if err != nil {
		helpers.WriteJSON(w, http.StatusUnauthorized, helpers.Envelope{"msg": err})
		return
	}
	tokenString, err := helpers.GenerateTokenString(r) // Generate token string
	if err != nil {
		helpers.WriteJSON(w, http.StatusUnauthorized, helpers.Envelope{"msg": err})
		return
	}

	claims, err := helpers.ValidClaim(&models.Claims{}, tokenString, jwtKey)
	if err != nil {
		helpers.WriteJSON(w, http.StatusUnauthorized, helpers.Envelope{"msg": err})
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

// GET -> /api/v1/donations/donation/{id}
func GetDonationByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": "Invalid id"})
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	donation, err := donation.GetDonationByID(id)
	if err != nil {
		helpers.WriteJSON(w, http.StatusNotFound, helpers.Envelope{"msg": "Not found"})
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, donation)
}

// GET -> /api/v1/donations/donations
func GetDonations(w http.ResponseWriter, r *http.Request) {
	allDonations, err := donation.GetAllDonations()
	if err != nil {
		helpers.WriteJSON(w, http.StatusNotFound, helpers.Envelope{"msg": "No donations found"})
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}
	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"donations": allDonations})
}

// Update donations
// PUT -> /api/v1/nodations/donate/{id}
func UpdateDonation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": "Invalid id"})
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&donationData)
	if err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": "Error while decoding the payload"})
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

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

	userID, _ := strconv.Atoi(claims.ID)
	msg, err := donation.UpdateDonation(id, userID, donationData)
	if err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": "Error while updating"})
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"msg": msg})
}

// Update donations
// DELETE -> /api/v1/nodations/donation/{id}
func DeleteDonation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, helpers.Envelope{"msg": "Invalid id"})
		helpers.MessageLogs.ErrorLog.Println(err)
		return
	}

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

	userID, _ := strconv.Atoi(claims.ID)
	msg, err := donation.DeleteDonation(id, userID)
	if err != nil {
		helpers.WriteJSON(w, http.StatusInternalServerError, helpers.Envelope{"msg": err})
	}

	helpers.WriteJSON(w, http.StatusOK, helpers.Envelope{"msg": msg})
	return
}
