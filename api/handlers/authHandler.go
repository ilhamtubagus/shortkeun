package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
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

// swagger:route POST /auth/signin auth signIn
//
//	Sign in (default)
//
//	Sign in with email and password
//
//  Security: None
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
func (a AuthHandler) SignIn(c echo.Context) error {
	var credential dto.SignInRequestDefaultBody
	err := c.Bind(&credential)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("failed to parse request body"))
	}
	//dto validation
	if err := c.Validate(credential); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity,
			dto.NewValidationError("validation failed", lib.MapError(err)))
	}
	usr, err := a.userRepository.FindUserByEmail(credential.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("unexpected database error"))
	}
	if usr == nil || usr.Password == "" {
		return echo.NewHTTPError(http.StatusNotFound, dto.NewDefaultResponse("user was not found"))
	}
	hasher := lib.NewBcryptHash()
	if err := hasher.CompareHash(credential.Password, usr.Password); err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("password does not match"))
	}
	hour, _ := strconv.Atoi(os.Getenv("TOKEN_EXP"))
	claims := entities.Claims{
		Role:   usr.Role,
		Email:  usr.Email,
		Status: usr.Status,
		StandardClaims: jwt.StandardClaims{
			//token expires within x hours
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(hour)).Unix(),
			Subject:   usr.ID.String(),
		}}
	token, err := claims.GenerateJwt()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("unexpected server error"))
	}
	return c.JSON(200, &dto.SignInResponseBody{Message: "signin succeeded", Token: *token})
}

//	swagger:route POST /auth/signin/google auth googleSignIn
//
//	Sign in with google account
//
//	Sign in with google account.
//	If user has not registered then registration process will be performed.
//
//	Consumes:
// 	- application/json
// 	Produces:
// 	- application/json
// 	Responses:
// 	- 422: validationError
//	- 200: signInResponse
//	- 404: defaultResponse
//	- 400: defaultResponse
//	- 500: defaultResponse
func (a AuthHandler) GoogleSignIn(c echo.Context) error {
	var credential dto.GoogleSignInRequestBody
	err := c.Bind(&credential)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("failed to parse request body"))
	}
	//dto validation
	if err := c.Validate(credential); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity,
			dto.NewValidationError("validation failed", lib.MapError(err)))
	}
	// decode and verify id token credential
	googleTokenInfo, err := lib.VerifyToken(credential.Credential)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	usr, err := a.userRepository.FindUserByEmail(googleTokenInfo.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("unexpected database error"))
	}
	if usr == nil {
		//insert new user into database
		usr = &entities.User{Name: googleTokenInfo.Name, Email: googleTokenInfo.Email, Sub: googleTokenInfo.Sub, Status: entities.StatusActive, Role: entities.RoleMember}
		err := a.userRepository.CreateUser(usr)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("unexpected database error"))
		}
	}
	//create our own jwt and send back to client
	hour, _ := strconv.Atoi(os.Getenv("TOKEN_EXP"))
	claims := entities.Claims{
		Role:   usr.Role,
		Email:  usr.Email,
		Status: usr.Status,
		StandardClaims: jwt.StandardClaims{
			//token expires within x hours
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(hour)).Unix(),
			Subject:   usr.ID.String(),
		}}
	token, err := claims.GenerateJwt()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dto.NewDefaultResponse("unexpected server error"))
	}
	return c.JSON(200, &dto.SignInResponseBody{Message: "signin succeeded", Token: *token})
}

//	swagger:route POST /auth/signin/register auth register
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
func (a AuthHandler) Register(c echo.Context) error {
	//dto binding
	registrant := new(dto.RegistrationRequestBody)
	if err := c.Bind(&registrant); err != nil {
		c.Echo().Logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, dto.NewDefaultResponse("failed to parse request body"))
	}
	//dto validation
	if err := c.Validate(registrant); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity,
			dto.NewValidationError("validation failed", lib.MapError(err)))
	}
	// domain validation
	// email must be unique for each users
	if user, err := a.userRepository.FindUserByEmail(registrant.Email); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else if user != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity,
			dto.NewValidationError("validation failed", &[]lib.ValidationError{
				{Field: "email", Message: "email has been registered"}}))
	}

	// perform password hashing
	hasher := lib.NewBcryptHash()
	hashedPassword, err := hasher.MakeHash(registrant.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	//generate activation code
	activationCode := lib.RandString(5)
	now := time.Now()
	//create user struct
	user := &entities.User{
		Name:     registrant.Name,
		Email:    registrant.Email,
		Password: *hashedPassword,
		Role:     entities.RoleMember,
		Status:   entities.StatusInactive,
		ActivationCode: &entities.ActivationCode{
			Code:     activationCode,
			IssuedAt: now,
			ExpireAt: now.Add(time.Minute * 5)},
	}
	err = a.userRepository.CreateUser(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	//send email registration with activation code
	ipE := echo.ExtractIPDirect()
	pathToTemplate, _ := filepath.Abs("./email/template/registrationMail.html")
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
	return c.JSON(http.StatusCreated, dto.NewDefaultResponse("registration succeeded"))
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
