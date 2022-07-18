package user

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ilhamtubagus/urlShortener/dto"
	"github.com/ilhamtubagus/urlShortener/lib"
	"github.com/labstack/echo/v4"
)

type UserService interface {
	FindUserByEmail(email string) (*User, error)
	SaveUser(user *User) error
	UpdateUser(user *User) error
	ActivateAccount(email, activationCode string) (*User, error)
	Register(user *User) (*User, error)
	RequestActivationCode(user *User) (*User, error)
}
type userService struct {
	userRepository UserRepository
	hash           lib.Hash
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

func (service userService) Register(user *User) (*User, error) {
	// domain validation
	// email must be unique for each users
	if user, err := service.userRepository.FindUserByEmail(user.Email); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else if user != nil {
		return nil, echo.NewHTTPError(http.StatusUnprocessableEntity,
			dto.NewValidationError("validation failed", &[]lib.ValidationError{
				{Field: "email", Message: "email has been registered"}}, http.StatusBadRequest))
	}
	// perform password hashing
	hashedPassword, err := service.hash.MakeHash(user.Password)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	user.Password = *hashedPassword
	//generate activation code
	activationCodeExpiry, _ := strconv.Atoi(os.Getenv("ACTIVATION_CODE_EXPIRY_IN_MINUTES"))
	activationCode := lib.RandString(5)
	now := time.Now()
	user.ActivationCode = &ActivationCode{
		Code:     activationCode,
		IssuedAt: now,
		ExpireAt: now.Add(time.Minute * time.Duration(activationCodeExpiry)),
	}
	err = service.userRepository.SaveUser(user)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return user, nil
}

func (service userService) ActivateAccount(email, activationCode string) (*User, error) {
	user, err := service.userRepository.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (authenticationService userService) RequestActivationCode(user *User) (*User, error) {
	user, err := authenticationService.userRepository.FindUserByEmail(user.Email)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if user == nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, dto.NewDefaultResponse("user with this email address was not found", http.StatusNotFound))
	}
	if user.Status == SUSPENDED {
		return nil, echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("user with this email address is suspended, please contact administrator for further information", http.StatusUnprocessableEntity))
	}
	if user.Status == ACTIVE {
		return nil, echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("user with this email address has already been activated", http.StatusBadRequest))
	}
	if user.ActivationCode != nil {
		if user.ActivationCode.ExpireAt.In(time.Now().Location()).After(time.Now()) {
			return nil, echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("the previous activation code has not been expired", http.StatusBadRequest))
		}
	}
	// issue new activation code
	now := time.Now()
	activationCode := lib.RandString(5)
	user.ActivationCode = &ActivationCode{
		Code:     activationCode,
		IssuedAt: now,
		ExpireAt: now.Add(5 * time.Minute),
	}
	err = authenticationService.userRepository.UpdateUser(user)
	if err != nil {
		fmt.Println(err)
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return user, nil
}
func NewUserService(userRepository UserRepository, hash lib.Hash) UserService {
	return userService{
		userRepository: userRepository,
		hash:           hash,
	}
}
