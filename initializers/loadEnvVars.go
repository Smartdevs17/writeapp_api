package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVars() {
	env := os.Getenv("GO_ENV")
	if env == "development" {
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	log.Printf("Environment: %s", env)
}
