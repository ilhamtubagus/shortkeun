package user

import (
	"github.com/ilhamtubagus/urlShortener/lib"
)

type UserService interface {
	FindUserByEmail(email string) (*User, error)
	SaveUser(user *User) error
	UpdateUser(user *User) error
	ActivateAccount(activationCode string) (*User, error)
}
type userService struct {
	hash           lib.Hash
	userRepository UserRepository
}

func (service userService) SaveUser(user *User) error {
	return service.userRepository.SaveUser(user)
}

func (service userService) FindUserByEmail(email string) (*User, error) {
	return service.userRepository.FindUserByEmail(email)
}

func (service userService) UpdateUser(user *User) error {
	return service.userRepository.UpdateUser(user)
}
func (service userService) ActivateAccount(email, activationCode string) (*User, error) {
	user, err := service.userRepository.FindUserByEmail(email)
}
