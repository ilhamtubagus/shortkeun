package service

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ilhamtubagus/urlShortener/domain/constant"
	"github.com/ilhamtubagus/urlShortener/domain/entity"
	"github.com/ilhamtubagus/urlShortener/domain/repository"
	"github.com/ilhamtubagus/urlShortener/interface/dto"
	"github.com/ilhamtubagus/urlShortener/utils"
	"github.com/labstack/echo/v4"
)

type UserService interface {
	FindUserByEmail(email string) (*entity.User, error)
	SaveUser(user *entity.User) error
	UpdateUser(user *entity.User) error
	ActivateAccount(user *entity.User, activationCode string) (*entity.User, error)
	Register(user *entity.User) (*entity.User, error)
	RequestActivationCode(user *entity.User) (*entity.User, error)
}
type userService struct {
	userRepository repository.UserRepository
	hash           utils.Hash
}

func (u userService) SaveUser(user *entity.User) error {
	return u.userRepository.CreateUser(user)
}

func (u userService) FindUserByEmail(email string) (*entity.User, error) {
	return u.userRepository.FindUserByEmail(email)
}

func (u userService) UpdateUser(user *entity.User) error {
	return u.userRepository.UpdateUser(user)
}

func (u userService) Register(user *entity.User) (*entity.User, error) {
	// domain validation
	// email must be unique for each users
	if user, err := u.userRepository.FindUserByEmail(user.Email); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else if user != nil {
		return nil, echo.NewHTTPError(http.StatusUnprocessableEntity,
			dto.NewValidationError("validation failed", &[]utils.ValidationError{
				{Field: "email", Message: "email has been registered"}}, http.StatusBadRequest))
	}
	// perform password hashing
	hashedPassword, err := u.hash.MakeHash(user.Password)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	user.Password = *hashedPassword
	//generate activation code
	activationCodeExpiry, _ := strconv.Atoi(os.Getenv("ACTIVATION_CODE_EXPIRY_IN_MINUTES"))
	activationCode := utils.RandString(5)
	now := time.Now()
	user.ActivationCode = &entity.ActivationCode{
		Code:     activationCode,
		IssuedAt: now,
		ExpireAt: now.Add(time.Minute * time.Duration(activationCodeExpiry)),
	}
	err = u.userRepository.CreateUser(user)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return user, nil
}

func (u userService) ActivateAccount(user *entity.User, activationCode string) (*entity.User, error) {
	searchedUser, err := u.userRepository.FindUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}
	if searchedUser == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("user not found", http.StatusBadRequest))
	}
	if searchedUser.Status != constant.SUSPENDED && searchedUser.Status != constant.INACTIVE {
		return nil, echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("current user status must be INACTIVE or SUSPENDED", http.StatusBadRequest))
	}
	if searchedUser.ActivationCode == nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("activation code was not found", http.StatusBadRequest))
	}
	if time.Now().After(searchedUser.ActivationCode.ExpireAt) {
		return nil, echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("activation code has been expired", http.StatusBadRequest))
	}
	searchedUser.Status = constant.ACTIVE
	err = u.userRepository.UpdateUser(searchedUser)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return searchedUser, nil
}

func (u userService) RequestActivationCode(user *entity.User) (*entity.User, error) {
	user, err := u.userRepository.FindUserByEmail(user.Email)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if user == nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, dto.NewDefaultResponse("user with this email address was not found", http.StatusNotFound))
	}
	if user.Status == constant.SUSPENDED {
		return nil, echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("user with this email address is suspended, please contact administrator for further information", http.StatusUnprocessableEntity))
	}
	if user.Status == constant.ACTIVE {
		return nil, echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("user with this email address has already been activated", http.StatusBadRequest))
	}
	if user.ActivationCode != nil {
		if time.Now().Before(user.ActivationCode.ExpireAt) {
			return nil, echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("the previous activation code has not been expired", http.StatusBadRequest))
		}
	}
	// issue new activation code
	now := time.Now()
	activationCode := utils.RandString(5)
	user.ActivationCode = &entity.ActivationCode{
		Code:     activationCode,
		IssuedAt: now,
		ExpireAt: now.Add(5 * time.Minute),
	}
	err = u.userRepository.UpdateUser(user)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return user, nil
}
func NewUserService(userRepository repository.UserRepository, hash utils.Hash) UserService {
	return userService{
		userRepository: userRepository,
		hash:           hash,
	}
}
