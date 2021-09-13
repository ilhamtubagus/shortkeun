package lib

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv(path string) {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatal("Error while loading .env " + err.Error())
	}
}
