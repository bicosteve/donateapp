package helpers

import (
	models2 "donateapp/pkg/models"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Envelope map[string]interface{}

type Message struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

var MessageLogs = &Message{
	InfoLog:  infoLog,
	ErrorLog: errorLog,
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	r.Close = true

	decode := json.NewDecoder(r.Body)

	err := decode.Decode(&data)
	if err != nil {
		return err
	}

	err = decode.Decode(&struct{}{})
	// Should not accept json like {}{} but accept this {{}}

	if err != io.EOF {
		return errors.New("invalid JSON value.")
	}

	return nil
}

func WriteJSON(
	w http.ResponseWriter, status int, data interface{}, headers ...http.Header,
) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload models2.JsonResponse
	payload.Error = true
	payload.Message = err.Error()

	err = WriteJSON(w, statusCode, payload)
	if err != nil {
		return err
	}

	return nil
}

func IsValidEmail(email string) bool {
	return govalidator.IsEmail(email)
}

func CheckPhoneNumber(phoneNumber string) bool {
	phoneNumber = strings.TrimSpace(phoneNumber)
	if phoneNumber == "" {
		return false
	}
	if len(phoneNumber) < 10 {
		return false
	}
	return true
}

func ValidatePassword(password string, confirmPassword string) bool {
	password = strings.TrimSpace(password)
	confirmPassword = strings.TrimSpace(confirmPassword)
	if password == "" {
		return false
	}
	if confirmPassword == "" {
		return false
	}
	if password != confirmPassword {
		return false
	}
	return true
}

func LoadJWTKEY() (string, error) {
	path, err := filepath.Abs(".env")
	if err != nil {
		return "", err
	}
	err = godotenv.Load(filepath.Join(path))
	if err != nil {
		return "", err
	}
	jwtKey := os.Getenv("JWTSECRET")
	return jwtKey, nil
}

func GenerateTokenString(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")

	if err != nil {
		if err == http.ErrNoCookie {
			return "", err
		}
		return "", err
	}

	tokenString := cookie.Value
	return tokenString, nil
}

func ValidClaim(
	claims *models2.Claims, tokenString string, jwtKey string,
) (*models2.Claims, error) {
	tkn, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Fatal(err)
			return nil, err
		}
		log.Fatal(err)
		return nil, err
	}

	if !tkn.Valid {
		log.Fatal(err)
		return nil, err
	}
	return claims, nil
}

func ValidateDonationPayload(donation models2.DonationBody) bool {
	name := strings.TrimSpace(donation.Name)
	photo := strings.TrimSpace(donation.Photo)
	location := strings.TrimSpace(donation.Location)

	if name == "" {
		return false
	}

	if photo == "" {
		return false
	}

	if location == "" {
		return false
	}

	return true
}
