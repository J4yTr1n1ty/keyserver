package boot

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Environment EnvHandler = EnvHandler{
	requiredEnvVars: []string{
		"PORT",
	},
}

type EnvHandler struct {
	requiredEnvVars []string
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	for _, requiredEnvVar := range Environment.requiredEnvVars {
		if os.Getenv(requiredEnvVar) == "" {
			panic("Required Environment variable not set: " + requiredEnvVar)
		}
	}
}

func (EnvHandler) GetEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic("Environment variable not set: " + key)
}

func (EnvHandler) GetEnvBool(key string) bool {
	if value, ok := os.LookupEnv(key); ok {
		return value == "true"
	}
	panic("Environment variable not set: " + key)
}
