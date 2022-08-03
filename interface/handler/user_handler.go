package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/ilhamtubagus/urlShortener/domain/constant"
	"github.com/ilhamtubagus/urlShortener/domain/entity"
	"github.com/ilhamtubagus/urlShortener/domain/service"
	"github.com/ilhamtubagus/urlShortener/email"
	"github.com/ilhamtubagus/urlShortener/interface/dto"
	"github.com/ilhamtubagus/urlShortener/utils"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
}

// swagger:route PATCH /user/status user accountActivation
//
//	Account activation
//
//	Activate user's account with activation code sent via email. User's status will change to "ACTIVE".
//
//  Security:
// 	- Bearer-Token:
//	Consumes:
// 	 application/json
// 	Produces:
// 	 application/json
// 	Responses:
// 	- 422: validationError
//	- 200: signInResponse
//	- 404: defaultResponse
//	- 400: defaultResponse
//	- 500: defaultResponse
func (u UserHandler) ActivateAccount(c echo.Context) error {
	accountActivationRequest := new(dto.AccountActivationRequestBody)
	if err := c.Bind(&accountActivationRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("failed to parse request body", http.StatusBadRequest))
	}
	// dto validation
	if err := c.Validate(accountActivationRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			dto.NewValidationError("Bad Request", utils.MapError(err), http.StatusUnprocessableEntity))
	}
	user, err := u.userService.ActivateAccount(&entity.User{
		Email: accountActivationRequest.Email,
	}, accountActivationRequest.ActivationCode)
	if err != nil {
		return err
	}
	userResponseDto := user.ConvertToResponseDto()
	return c.JSON(http.StatusOK, &userResponseDto)
}

//	swagger:route POST /users register
//
//  Register new account
//
//	Register new account with email and password.
//	User will be given a code for account activation via email after registration has been performed.
//
//	Consumes:
// 	- application/json
// 	Produces:
// 	- application/json
// 	Responses:
// 	- 422: validationError
//	- 200: defaultResponse
//	- 404: defaultResponse
//	- 400: defaultResponse
//	- 500: defaultResponse
func (u UserHandler) Register(c echo.Context) error {
	registrant := new(dto.RegistrationRequestBody)
	if err := c.Bind(&registrant); err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("failed to parse request body", http.StatusBadRequest))
	}

	if err := c.Validate(registrant); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity,
			dto.NewValidationError("validation failed", utils.MapError(err), http.StatusUnprocessableEntity))
	}
	user := &entity.User{
		Name:     registrant.Name,
		Email:    registrant.Email,
		Password: registrant.Password,
		Role:     constant.MEMBER,
		Status:   constant.INACTIVE,
	}
	user, err := u.userService.Register(user)
	if err != nil {
		return err
	}
	//send email registration with activation code
	ipE := echo.ExtractIPDirect()
	now := time.Now()
	pathToTemplate, _ := filepath.Abs("./email/template/registrationMail.html")
	attachment, _ := filepath.Abs("./logo.png")
	emailBody := email.RegistrationEmailBody{
		UserAgent: c.Request().UserAgent(),
		IP:        ipE(c.Request()),
		DateTime:  now.Format("Monday, 02-Jan-06 15:04:05 MST"),
		Code:      user.ActivationCode.Code,
		ExpireAt:  user.ActivationCode.ExpireAt.Format("Monday, 02-Jan-06 15:04:05 MST"),
	}
	// asynchronously send email registration
	go func() {
		err := utils.SendHTMLMail([]string{user.Email}, "Activate Your Account", emailBody, pathToTemplate, []string{attachment})
		if err != nil {
			c.Logger().Error(fmt.Sprintf("failed to send email registration to %s", user.Email))
		}
	}()
	return c.JSON(http.StatusCreated, &user)
}

//	swagger:route POST /users/activation-code auth getActivationCode
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
func (u UserHandler) RequestActivationCode(c echo.Context) error {
	newActivationCodeRequest := new(dto.NewActivationCodeRequestBody)
	if err := c.Bind(&newActivationCodeRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("failed to parse request body", http.StatusBadRequest))
	}
	// dto validation
	if err := c.Validate(newActivationCodeRequest); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity,
			dto.NewValidationError("validation failed", utils.MapError(err), http.StatusUnprocessableEntity))
	}

	user := &entity.User{
		Email: newActivationCodeRequest.Email,
	}
	user, err := u.userService.RequestActivationCode(user)
	if err != nil {
		return err
	}
	//send email registration with activation code
	ipE := echo.ExtractIPDirect()
	now := time.Now()
	pathToTemplate, _ := filepath.Abs("./email/template/activationMail.html")
	attachment, _ := filepath.Abs("./logo.png")
	emailBody := email.RegistrationEmailBody{
		UserAgent: c.Request().UserAgent(),
		IP:        ipE(c.Request()),
		DateTime:  now.Format("Monday, 02-Jan-06 15:04:05 MST"),
		Code:      user.ActivationCode.Code,
		ExpireAt:  user.ActivationCode.ExpireAt.Format("Monday, 02-Jan-06 15:04:05 MST"),
	}
	// asynchronously send email registration
	go func() {
		err := utils.SendHTMLMail([]string{user.Email}, "Activate Your Account", emailBody, pathToTemplate, []string{attachment})
		if err != nil {
			c.Logger().Error(fmt.Sprintf("failed to send email registration to %s", user.Email))
		}
	}()
	return c.JSON(http.StatusCreated, dto.NewDefaultResponse("activation code sent", http.StatusOK))
}

func NewUserHandler(userService service.UserService) UserHandler {
	return UserHandler{userService: userService}
}
