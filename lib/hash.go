package lib

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Hasher interface {
	MakeHash(plain string) (*string, error)
	CompareHash(plain, hashToCompare string) error
}
type BcryptHasher struct {
}

func (b BcryptHasher) MakeHash(plain string) (*string, error) {
	if plain == "" {
		return nil, errors.New("plain string must not be an empty string")
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	hashedTxt := string(bytes)
	return &hashedTxt, err
}
func (b BcryptHasher) CompareHash(plain, hashToCompare string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashToCompare), []byte(plain))
}
func NewBcryptHasher() Hasher {
	return &BcryptHasher{}
}
