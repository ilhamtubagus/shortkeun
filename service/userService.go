package service

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ilhamtubagus/urlShortener/dto"
	"github.com/ilhamtubagus/urlShortener/entities"
	"github.com/ilhamtubagus/urlShortener/repository"
)

type UserService struct {
	userRepository repository.UserRepository
}

func (u UserService) ActivateAccount(userId, activationCode string) *dto.ApiError {
	user, err := u.userRepository.FindUserById(userId)
	if err != nil {
		fmt.Println(err)
		return dto.NewApiError(http.StatusInternalServerError, errors.New("unexpected database error"))
	}
	if user == nil {
		return dto.NewApiError(http.StatusNotFound, errors.New("user was not found"))
	}
	if user.ActivationCode == nil {
		return dto.NewApiError(http.StatusNotFound, errors.New("activation code is not present"))
	}
	if user.ActivationCode.ExpireAt.In(time.Now().Location()).Before(time.Now()) {
		return dto.NewApiError(http.StatusBadRequest, errors.New("the previous activation code has been expired, please request the new one"))
	}
	if user.ActivationCode.Code != activationCode {
		return dto.NewApiError(http.StatusUnprocessableEntity, errors.New("activation code is incorrect"))
	}
	user.Status = entities.StatusActive
	user.ActivationCode = nil
	u.userRepository.UpdateUser(user)
	return nil
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}
