package main

import (
	"github.com/gabbottron/messenger-api/api"
	"github.com/gabbottron/messenger-api/src/datastore"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func GetAppEnv() error {
	env := os.Getenv("APP_ENV")

	if len(env) == 0 {
		err := godotenv.Load()
		if err != nil {
			log.Panic("Error loading .env file")
			return err
		}
	}

	return nil
}

func main() {
	log.Println("Starting the messenger API...")

	// Make sure the application can load env
	err := GetAppEnv()
	if err != nil {
		log.Panic("Error loading APP ENV!")
	}

	// Connect to the database
	err = datastore.InitDB()
	if err != nil {
		log.Panic("Error connecting to the DB")
	}

	// Initialize the router and run it
	router := api.InitRouter()
	router.Run()
}
