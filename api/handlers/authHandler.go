package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ilhamtubagus/urlShortener/dto"
	"github.com/ilhamtubagus/urlShortener/email"
	"github.com/ilhamtubagus/urlShortener/entities"
	"github.com/ilhamtubagus/urlShortener/lib"
	"github.com/ilhamtubagus/urlShortener/repository"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	userRepository repository.UserRepository
}

//	swagger:route POST /auth/signin/activation-code auth getActivationCode
//
//	Get new activation code
//
//	Get new activation code for account activation purpose if the previous activation code has been expired.
//
//	Consumes:
// 	- application/json
// 	Produces:
// 	- application/json
// 	Responses:
// 	- 422: validationError
//	- 201: defaultResponse
//	- 404: defaultResponse
//	- 400: defaultResponse
//	- 500: defaultResponse
func (ah AuthHandler) RequestActivationCode(c echo.Context) error {
	requestCodeAct := new(dto.ActivationCodeRequestBody)
	if err := c.Bind(&requestCodeAct); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("failed to parse request body"))
	}
	// dto validation
	if err := c.Validate(requestCodeAct); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity,
			dto.NewValidationError("validation failed", lib.MapError(err)))
	}
	user, err := ah.userRepository.FindUserByEmail(requestCodeAct.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound, dto.NewDefaultResponse("user with this email address was not found"))
	}
	if user.Status == entities.StatusSuspended {
		return echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("user with this email address is suspended, please contact administrator for further information"))
	}
	if user.Status == entities.StatusActive {
		return echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("user with this email address has already been activated"))
	}
	if user.ActivationCode != nil {
		if user.ActivationCode.ExpireAt.In(time.Now().Location()).After(time.Now()) {
			return echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("the previous activation code has not been expired"))
		}
	}

	// issue new activation code
	now := time.Now()
	activationCode := lib.RandString(5)
	user.ActivationCode = &entities.ActivationCode{
		Code:     activationCode,
		IssuedAt: now,
		ExpireAt: now.Add(5 * time.Minute),
	}
	err = ah.userRepository.UpdateUser(user)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	//send email registration with activation code
	ipE := echo.ExtractIPDirect()
	pathToTemplate, _ := filepath.Abs("./email/template/activationMail.html")
	attachment, _ := filepath.Abs("./logo.png")
	emailBody := email.RegistrationMailBody{
		UserAgent: c.Request().UserAgent(),
		IP:        ipE(c.Request()),
		DateTime:  now.Format("Monday, 02-Jan-06 15:04:05 MST"),
		Code:      user.ActivationCode.Code,
		ExpireAt:  user.ActivationCode.ExpireAt.Format("Monday, 02-Jan-06 15:04:05 MST"),
	}
	// asynchronously send email registration
	go func() {
		err := lib.SendHTMLMail([]string{user.Email}, "Activate Your Account", emailBody, pathToTemplate, []string{attachment})
		if err != nil {
			c.Logger().Error(fmt.Sprintf("failed to send email registration to %s", user.Email))
		}
	}()
	return c.JSON(http.StatusCreated, dto.NewDefaultResponse("activation code sent"))
}

// todo refresh token

func NewAuthHandler(userRepository repository.UserRepository) AuthHandler {
	return AuthHandler{userRepository: userRepository}
}
