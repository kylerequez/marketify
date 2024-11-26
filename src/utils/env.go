package utils

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	log.Println("Loading .env file...")
	err := godotenv.Load()
	if err != nil {
		return err
	}
	log.Println("Successfully loaded .env file")
	return nil
}

func GetEnv(key string) (*string, error) {
	value := os.Getenv(key)
	if value == "" {
		return nil, errors.New(fmt.Sprintf("%s is not found", key))
	}
	return &value, nil
}
