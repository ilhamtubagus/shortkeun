package lib

import (
	"errors"

	"github.com/joho/godotenv"
)

func LoadEnv(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return errors.New("Error while loading .env " + err.Error())
	}
	return nil
}
