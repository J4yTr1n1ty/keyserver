package boot

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Environment EnvHandler = EnvHandler{
	requiredEnvVars: []string{
		"PORT",
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASS",
		"DB_NAME",
	},
}

type EnvHandler struct {
	requiredEnvVars []string
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Unable to load env file, continuing without.")
	}

	sucess := true
	missingEnvVars := make([]string, 0)
	for _, requiredEnvVar := range Environment.requiredEnvVars {
		if _, ok := os.LookupEnv(requiredEnvVar); !ok {
			sucess = false
			missingEnvVars = append(missingEnvVars, requiredEnvVar)
		}
	}

	if !sucess {
		panic("Missing environment variables: " + fmt.Sprint(missingEnvVars))
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
