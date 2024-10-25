package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type (
	Env struct {
		// DBConnectionString the connections string
		DBConnectionString string
		// ServerHost the host that the server will start on
		ServerHost string
		// ServerPort the port that server will start on
		ServerPort string
	}
)

// NewEnv creates a new instance of Env
// tries to load the env variables from .env
func NewEnv() *Env {
	path, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting path")
	}
	dotenvError := godotenv.Load(fmt.Sprintf("%s/.env", path))
	if dotenvError != nil {
		log.Printf("error loading .env file, ignoring dotenv")
	}

	return &Env{
		DBConnectionString: os.Getenv("DB_CONNECTION_STRING"),
		ServerHost:         os.Getenv("Server_Host"),
		ServerPort:         os.Getenv("Server_Port"),
	}
}
