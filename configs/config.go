package configs

import (
	"log"
	"syscall"

	"github.com/joho/godotenv"
)

func init() {
	loadEnv()
}

func loadEnv() {
	err := godotenv.Overload()
	if err != nil {
		log.Fatalf("error in loading environment variables")
	}
}

func GetEnvWithKey(key string, defaultVal string) string {
	val, isFound := syscall.Getenv(key)
	if !isFound {
		syscall.Setenv(key, defaultVal)
		return defaultVal
	}
	return val
}
