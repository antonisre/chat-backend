package main

import (
	"chat-backend/app/socket"
	"chat-backend/config"
	"chat-backend/routes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	APIRoutes := routes.Api{}

	// load env file
	var env string = "local"
	envVariable := os.Getenv("APP_ENV")

	if envVariable == "development" {
		env = "development"
	} else if envVariable == "production" {
		env = "production"
	}

	if err := godotenv.Load(".env." + env); err != nil {
		fmt.Println("Failed to load env")
		panic(err)
	}

	// Init the database
	config.Connect(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	// Init the routes
	APIRoutes.ServeRoutes()

	//init socket writer
	go socket.WriteSocketMessage()

	// Run the server
	fmt.Println("Connected To Database")
	fmt.Println("Server started port 8000")
	log.Fatal(http.ListenAndServe(":8000", APIRoutes.Router))
}
