package u

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

//Obtiene variables de entorno
func Get(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
