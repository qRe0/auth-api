package auth

import (
	"context"
	"fmt"
	"log"
	"math/rand"
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

func (a *AuthService) SignUp(ctx context.Context, user *models.User) error {
	if user.Name == "" || user.Phone == "" || user.Email == "" || user.Password == "" {
		return errors.Wrap(errs.ErrIncorrectData, "failed to validate registration data. incorrect user data")
	}

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return errors.Wrap(err, "failed to hash password")
	}

	user.Password = hashedPassword
	err = a.repo.CreateUser(ctx, user)
	if err != nil {
		return errors.Wrap(err, errs.ErrCreateUser.Error())
	}

	return nil
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

	token, err := a.NewSession(ctx, userDataFromDB.ID, a.cfg.SecretKey, 5*time.Minute)
	if err != nil {
		return models.Tokens{}, errors.Wrap(err, "failed to create new session")
	}

	return token, nil
}

func validatePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) NewSession(ctx context.Context, userID string, secretKey string, lifetime time.Duration) (models.Tokens, error) {
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

func newJWT(userID string, secretKey string, lifetime time.Duration) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(lifetime).Unix()

	log.Printf("claims: %+v\n", claims)

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
