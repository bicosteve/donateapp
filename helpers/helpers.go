package helpers

import (
	"donateapp/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
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

	decode := json.NewDecoder(r.Body)

	err := decode.Decode(data)

	if err != nil {
		return err
	}

	err = decode.Decode(&struct{}{})
	// Should not have {}{} but {{}}

	if err != nil {
		return errors.New("Invalid JSON object")
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
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

	var payload models.JsonResponse
	payload.Error = true
	payload.Message = err.Error()

	err = WriteJSON(w, statusCode, payload)
	if err != nil {
		return err
	}

	return nil
}

//// DB Helper
//func UserExists(db *sql.DB, user User) (bool, error) {
//	row := db.QueryRow("SELECT id FROM users WHERE email LIKE %?%", user.Email)
//	var id int64
//
//	err := row.Scan(&id)
//
//	if err == sql.ErrNoRows {
//		return false, nil
//	}
//
//	if err != nil {
//		return false, err
//	}
//
//	return true, nil
//}
