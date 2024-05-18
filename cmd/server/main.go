package main

import (
	"donateapp/pkg/db"
	"donateapp/pkg/models"
	"donateapp/pkg/routes"
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
	Config Config
	Models models.Models
}

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Cannot load .env file")

	}

	config := Config{
		Port: os.Getenv("PORT"),
	}

	dbHost := os.Getenv("DBHOST")
	dbUser := os.Getenv("DBUSER")
	dbPort := os.Getenv("DBPORT")
	dbPassword := os.Getenv("DBPASSWORD")
	dbName := os.Getenv("DBNAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	dbConnection, err := db.ConnectMysql(dsn)

	if err != nil {
		log.Fatal("Cannot connect to sql")
	}

	defer dbConnection.DB.Close()

	app := &Application{
		Config: config,
		Models: models.NewConnections(dbConnection.DB),
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}

}

func (app *Application) serve() error {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file for app configs")
	}

	port := os.Getenv("PORT")
	fmt.Printf("Listening to port %v ....\n", port)

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		// To add route handler
		Handler: routes.Routes(),
	}

	return server.ListenAndServe()
}
