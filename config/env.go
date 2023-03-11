package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

// Load environment variables and handle errors

func LoadEnv() {
	err := godotenv.Load()

	if err!=nil {
		fmt.Println("Failed to load env")
		log.Fatal(err)
	}
}
