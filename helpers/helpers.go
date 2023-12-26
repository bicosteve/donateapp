package helpers

import (
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
