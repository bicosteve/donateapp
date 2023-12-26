package appconfigs

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

	rootDir, err := os.Getwd()

	if err != nil {
		log.Fatal("Error getting current directory")
	}

	envPath := filepath.Join(rootDir, ".env")

	err = godotenv.Load(envPath)

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
