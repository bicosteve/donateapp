package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

type Application struct {
	//Add Configs and Models
	Config Config
	// Todo Add Models
}

func (app *Application) Serve() error {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatal("Error loading .env file for app configs")
	}

	port := os.Getenv("PORT")
	fmt.Printf("Listening to port %v ....\n", port)

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		// To add route handler
	}

	return server.ListenAndServe()
}

func main() {

	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatal("Cannot load .env file")

	}

	config := Config{
		Port: os.Getenv("PORT"),
	}

	// To add the db connection dsn -> host, user, port, password, dbname
	// Create dsn string with fmt.Sprintf

	app := &Application{
		Config: config,
		// To add Models
	}

	err = app.Serve()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(".env file loaded")

}
