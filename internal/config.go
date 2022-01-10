package internal

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// GetKey ReadKey  func to get env value
func GetKey(key string) string {

	//_, ok := os.LookupEnv("debug")
	//if ok {

	// load .env file so it's accessible from GetEnv
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	return os.Getenv(key)
	//}
}
