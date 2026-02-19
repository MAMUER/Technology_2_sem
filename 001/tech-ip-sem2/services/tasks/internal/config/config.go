package config

import (
	"os"
)

type Config struct {
	Port        string 
	AuthBaseURL string
}

func Load() Config {
	port := os.Getenv("TASKS_PORT")
	if port == "" {
		port = "8082"
	}

	authBaseURL := os.Getenv("AUTH_BASE_URL")
	if authBaseURL == "" {
		authBaseURL = "http://localhost:8081"
	}

	return Config{
		Port:        port,
		AuthBaseURL: authBaseURL,
	}
}
