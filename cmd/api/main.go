package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatal("Cannot load .env file")
	}

	fmt.Println(".env file loaded")

}