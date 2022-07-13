package user

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ilhamtubagus/urlShortener/dto"
	"github.com/ilhamtubagus/urlShortener/lib"
	"github.com/labstack/echo/v4"
)

type UserService interface {
	Register(user *User) (*User, error)
}
type userService struct {
	hash           lib.Hash
	userRepository UserRepository
}

func (service userService) Register(user *User) (*User, error) {
	// domain validation
	// email must be unique for each users
	if user, err := service.userRepository.FindUserByEmail(user.Email); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else if user != nil {
		return nil, echo.NewHTTPError(http.StatusUnprocessableEntity,
			dto.NewValidationError("validation failed", &[]lib.ValidationError{
				{Field: "email", Message: "email has been registered"}}))
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
	return user, nil
}
