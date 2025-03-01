package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func Init() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	if os.Getenv("DEBUG") == "true" {
		log.Print("Env initialized")

	}
}
