package lib

import "golang.org/x/crypto/bcrypt"

type Hasher interface {
	MakeHash(plain string) (string, error)
	CompareHash(plain, hashToCompare string) (bool, error)
}
type BcryptHasher struct {
}

func (b BcryptHasher) MakeHash(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	hashedTxt := string(bytes)
	return hashedTxt, err
}
func (b BcryptHasher) CompareHash(plain, hashToCompare string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashToCompare), []byte(plain))
	if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
func NewBcryptHasher() Hasher {
	return &BcryptHasher{}
}
