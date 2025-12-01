package config

import (
	"os"

	"github.com/joho/godotenv"
)

var mandatoryEnvVars = []string{
	"POSTGRES_URL",
	"BETTER_AUTH_URL",
}

func InitEnv(filenames ...string) {
	godotenv.Load(filenames...)

	checkEnv()
}

func checkEnv() {
	for _, envVar := range mandatoryEnvVars {
		if value, ok := os.LookupEnv(envVar); !ok {
			panic(envVar + " is not set in the environment variables")
		} else if value == "" {
			panic(envVar + " is set but empty")
		}
	}
}
