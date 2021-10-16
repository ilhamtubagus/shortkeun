package service

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/ilhamtubagus/urlShortener/dto"
	"github.com/ilhamtubagus/urlShortener/entities"
	"github.com/ilhamtubagus/urlShortener/lib"
	"github.com/ilhamtubagus/urlShortener/repository"
)

type AuthService struct {
	userRepository repository.UserRepository
	hasher         lib.Hasher
}

func generateToken(usr *entities.User) (*entities.Token, error) {
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
	return claims.GenerateJwt()
}

func (as AuthService) SignIn(credential *dto.SignInRequestDefaultBody) (*entities.Token, *dto.ApiError) {
	usr, err := as.userRepository.FindUserByEmail(credential.Email)
	if err != nil {
		return nil, dto.NewApiError(http.StatusInternalServerError, errors.New("unexpected database error"))
	}
	// if user is nill or user's password is empty (signed with google account) then return 404 error
	if usr == nil || usr.Password == "" {
		return nil, dto.NewApiError(http.StatusInternalServerError, errors.New("user was not found"))
	}
	// compare password
	if err := as.hasher.CompareHash(credential.Password, usr.Password); err != nil {
		return nil, dto.NewApiError(http.StatusUnauthorized, errors.New("password does not match"))
	}
	token, err := generateToken(usr)
	if err != nil {
		return nil, dto.NewApiError(http.StatusInternalServerError, errors.New("unexpected server error"))
	}
	return token, nil
}

func (as AuthService) GoogleSignIn(credential *dto.GoogleSignInRequestBody) (*entities.Token, *dto.ApiError) {
	// decode and verify id token credential
	googleTokenInfo, err := lib.VerifyToken(credential.Credential)
	if err != nil {
		return nil, dto.NewApiError(http.StatusUnauthorized, err)
	}
	usr, err := as.userRepository.FindUserByEmail(googleTokenInfo.Email)
	if err != nil {
		return nil, dto.NewApiError(http.StatusInternalServerError, errors.New("unexpected database error"))
	}
	if usr == nil {
		//insert new user into database
		usr = &entities.User{Name: googleTokenInfo.Name, Email: googleTokenInfo.Email, Sub: googleTokenInfo.Sub, Status: entities.StatusActive, Role: entities.RoleMember}
		err := as.userRepository.CreateUser(usr)
		if err != nil {
			return nil, dto.NewApiError(http.StatusInternalServerError, errors.New("unexpected database error"))
		}
	}
	//create our own jwt and send back to client
	token, err := generateToken(usr)
	if err != nil {
		return nil, dto.NewApiError(http.StatusInternalServerError, errors.New("unexpected server error"))
	}
	return token, nil
}

func (as AuthService) Register(registrant *dto.RegistrationRequestBody) (*entities.User, *dto.ApiError) {
	// domain validation
	// email must be unique for each users
	if user, err := as.userRepository.FindUserByEmail(registrant.Email); err != nil {
		return nil, dto.NewApiError(http.StatusInternalServerError, errors.New("unexpected database error"))
	} else if user != nil {
		return nil, dto.NewApiError(http.StatusUnprocessableEntity, errors.New("email has been registered"))
	}
	// perform password hashing
	hashedPassword, err := as.hasher.MakeHash(registrant.Password)
	if err != nil {
		return nil, dto.NewApiError(http.StatusInternalServerError, errors.New("unexpected server error"))
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
	err = as.userRepository.CreateUser(user)
	if err != nil {
		return nil, dto.NewApiError(http.StatusInternalServerError, errors.New("unexpected database error"))
	}
	return user, nil
}

func (as AuthService) RequestActivationCode(requestCodeAct *dto.ActivationCodeRequestBody) (*entities.User, *dto.ApiError) {
	user, err := as.userRepository.FindUserByEmail(requestCodeAct.Email)
	if err != nil {
		return nil, dto.NewApiError(http.StatusInternalServerError, errors.New("unexpected database error"))
	}
	if user == nil {
		return nil, dto.NewApiError(http.StatusNotFound, errors.New("user was not found"))
	}
	// user with status "SUSPENDED" can't request activation code to prevent its account being activated
	if user.Status == entities.StatusSuspended {
		return nil, dto.NewApiError(http.StatusBadRequest, errors.New("user suspended, can't request activation code"))
	}
	// user with status "ACTIVE" cant request activation code
	if user.Status == entities.StatusActive {
		return nil, dto.NewApiError(http.StatusBadRequest, errors.New("user's status is active, cant request new activation code"))
	}
	// if previous activation code exist check if it has expired
	if user.ActivationCode != nil {
		if user.ActivationCode.ExpireAt.In(time.Now().Location()).After(time.Now()) {
			return nil, dto.NewApiError(http.StatusBadRequest, errors.New("the previous activation code has not been expired"))
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
	err = as.userRepository.UpdateUser(user)
	if err != nil {
		return nil, dto.NewApiError(http.StatusInternalServerError, errors.New("unexpected database error"))
	}
	return user, nil
}

func NewAuthService(userRepository repository.UserRepository, hasher lib.Hasher) *AuthService {
	return &AuthService{userRepository: userRepository, hasher: hasher}
}
