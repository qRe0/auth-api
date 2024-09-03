package auth

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"github.com/qRe0/auth-api/configs"
	errs "github.com/qRe0/auth-api/internal/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/qRe0/auth-api/internal/models"
	repository "github.com/qRe0/auth-api/internal/repository/auth"
	tokenService "github.com/qRe0/auth-api/internal/service/token"
)

type AuthService struct {
	cfg       configs.JWTConfig
	repo      repository.AuthRepositoryInterface
	tokenServ tokenService.TokenServiceInterface
}

func NewAuthService(cfg configs.JWTConfig, repo repository.AuthRepositoryInterface, service tokenService.TokenServiceInterface) *AuthService {
	return &AuthService{
		cfg:       cfg,
		repo:      repo,
		tokenServ: service,
	}
}

func (a *AuthService) SignUp(ctx context.Context, user *models.User) (models.Tokens, error) {
	if user.Name == "" || user.Phone == "" || user.Email == "" || user.Password == "" {
		return models.Tokens{}, errors.Wrap(errs.ErrIncorrectData, "failed to validate registration data. incorrect user data")
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return models.Tokens{}, errors.Wrap(err, "failed to hash password")
	}

	user.Password = hashedPassword
	err = a.repo.CreateUser(ctx, user)
	if err != nil {
		return models.Tokens{}, errors.Wrap(err, errs.ErrCreateUser.Error())
	}

	userDataFromDB, err := a.repo.GetUserData(ctx, user.Phone)
	if err != nil {
		return models.Tokens{}, errors.Wrap(err, errs.ErrUserNotFound.Error())
	}

	id, err := strconv.Atoi(userDataFromDB.ID)
	if err != nil {
		return models.Tokens{}, errors.Wrap(errs.ErrIncorrectData, "incorrect user id")
	}
	tokens, err := a.NewSession(ctx, id, a.cfg.SecretKey, 5*time.Minute)

	return tokens, nil
}

func hashPassword(password string) (string, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	if err != nil {
		return "", err
	}

	return string(hashedPass), nil
}

func (a *AuthService) LogIn(ctx context.Context, user *models.User) (models.Tokens, error) {
	if user.Phone == "" || user.Password == "" {
		return models.Tokens{}, errors.Wrap(errs.ErrIncorrectData, "failed to validate login data. incorrect phone number or password")
	}

	userDataFromDB, err := a.repo.GetUserData(ctx, user.Phone)
	if err != nil {
		return models.Tokens{}, errs.ErrUserNotFound
	}

	err = validatePassword(userDataFromDB.Password, user.Password)
	if err != nil {
		return models.Tokens{}, errors.Wrap(errs.ErrIncorrectData, "wrong password")
	}

	userID, err := strconv.Atoi(userDataFromDB.ID)
	if err != nil {
		return models.Tokens{}, errors.Wrap(errs.ErrIncorrectData, "incorrect user id")
	}

	accessToken, err := newJWT(userID, a.cfg.SecretKey, 5*time.Minute)
	if err != nil {
		return models.Tokens{}, err
	}

	refreshTokenExistsResponse, err := a.tokenServ.RefreshTokenExists(ctx, userID)
	if err != nil {
		return models.Tokens{}, err
	}

	var refreshToken string
	if !refreshTokenExistsResponse.Exists {
		refreshToken, err = newRefreshToken()
		if err != nil {
			return models.Tokens{}, err
		}
	} else {
		refreshToken = refreshTokenExistsResponse.RefreshToken
	}

	tokens := models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokens, nil
}

func validatePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) NewSession(ctx context.Context, userID int, secretKey string, lifetime time.Duration) (models.Tokens, error) {
	accessToken, err := newJWT(userID, secretKey, lifetime)
	if err != nil {
		return models.Tokens{}, err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return models.Tokens{}, err
	}

	err = a.tokenServ.SaveToken(ctx, refreshToken, userID)
	if err != nil {
		return models.Tokens{}, err
	}

	return models.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, err
}

func newJWT(userID int, secretKey string, lifetime time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(lifetime).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", errs.ErrGenerateToken
	}
	return signedToken, nil
}

func newRefreshToken() (string, error) {
	token := make([]byte, 32)

	source := rand.NewSource(time.Now().Unix())
	randomizer := rand.New(source)

	_, err := randomizer.Read(token)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", token), nil
}

func (a *AuthService) Refresh(ctx context.Context, token string) (models.Tokens, error) {
	userID, err := a.tokenServ.GetUserIDByRefreshToken(ctx, token)
	if err != nil {
		return models.Tokens{}, err
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return models.Tokens{}, errors.Wrap(errs.ErrIncorrectData, "incorrect user id")
	}

	newToken, err := a.NewSession(ctx, id, a.cfg.SecretKey, 5*time.Minute)
	if err != nil {
		return models.Tokens{}, err
	}

	return newToken, nil
}

func (a *AuthService) RevokeTokens(ctx context.Context, user *models.User) error {
	if user.Phone == "" || user.Email == "" || user.Password == "" || user.Name == "" {
		return errors.Wrap(errs.ErrIncorrectData, "incorrect input number")
	}

	userDataFromDB, err := a.repo.GetUserData(ctx, user.Phone)
	if err != nil {
		return errs.ErrUserNotFound
	}

	err = validatePassword(userDataFromDB.Password, user.Password)
	if err != nil {
		return errors.Wrap(errs.ErrIncorrectData, "wrong password")
	}

	id, err := strconv.Atoi(userDataFromDB.ID)
	if err != nil {
		return errors.Wrap(errs.ErrIncorrectData, "incorrect user id")
	}

	err = a.tokenServ.DeleteToken(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) ValidateToken(token string, cfg string) (string, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg), nil
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to parse token")
	}

	userID := fmt.Sprintf("%v", claims["sub"])
	return userID, nil
}

func (a *AuthService) TokenBlacklisted(ctx context.Context, token string) (bool, error) {
	if token == "" {
		return false, errors.Wrap(errs.ErrIncorrectData, "empty token")
	} else {
		return a.tokenServ.TokenBlacklisted(ctx, token)
	}
}

func (a *AuthService) LogOut(ctx context.Context, token string) error {
	if token == "" {
		return errors.Wrap(errs.ErrIncorrectData, "empty token")
	}

	splitedToken := strings.Split(token, " ")
	if len(splitedToken) != 2 {
		return errors.Wrap(errs.ErrIncorrectData, "incorrect token format")
	}
	token = splitedToken[1]

	userID, err := a.ValidateToken(token, a.cfg.SecretKey)
	if err != nil {
		return errors.Wrap(err, "logout failed")
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return errors.Wrap(errs.ErrIncorrectData, "incorrect user id")
	}

	err = a.tokenServ.DeleteToken(ctx, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete token")
	}

	err = a.tokenServ.BlacklistToken(ctx, token)
	if err != nil {
		return errors.Wrap(err, "failed to blacklist token")
	}

	return nil
}
