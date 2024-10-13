package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"terminal-resume.jayash.space/api"
	"terminal-resume.jayash.space/models"
)

func main() {

	if err := godotenv.Load(); err !=nil {
		log.Fatal("Error loading .env file")
	}

	config := models.Config{
        Host:     os.Getenv("HOST"),
        Port:     os.Getenv("PORT"),
        User:     os.Getenv("PGUSER"),
        Password: os.Getenv("PASSWORD"),
        DBName:   os.Getenv("DBNAME"),
        SSLMode:  os.Getenv("SSLMODE"),
    }

	models.InitDB(config)

	apiServer := api.NewApiServer(":8080")
	apiServer.Start()
}
