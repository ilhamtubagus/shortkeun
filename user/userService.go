package user

import (
	"github.com/ilhamtubagus/urlShortener/lib"
)

type UserService interface {
	FindUserByEmail(email string) (*User, error)
	SaveUser(user *User) error
	UpdateUser(user *User) error
}
type userService struct {
	hash           lib.Hash
	userRepository UserRepository
}

func (service userService) SaveUser(user *User) error {
	err := service.userRepository.SaveUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (service userService) FindUserByEmail(email string) (*User, error) {
	user, err := service.userRepository.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service userService) UpdateUser(user *User) error {
	err := service.userRepository.UpdateUser(user)
	if err != nil {
		return err
	}
	return nil
}
