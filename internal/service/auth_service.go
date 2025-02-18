package service

import (
	"time"

	"com.mx/crud/config"
	"com.mx/crud/internal/models"
	"com.mx/crud/internal/repository"
	"com.mx/crud/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenExpiryDuration  = time.Hour * 1      // 1 hour
	refreshTokenExpiryDuration = time.Hour * 24 * 7 // 7 days
	EMAIL                      = "email"
	USERNAME                   = "username"
)

type AuthService interface {
	RegisterUser(*models.User) (*models.User, error)
	LoginUser(email, password string) (*models.User, string, string, error)
	RefreshAccessToken(refreshToken string) (string, string, error)
}

type authService struct {
	userRepo    repository.UserRepository
	userService UserService
}

func NewAuthService(userRepo repository.UserRepository, userService UserService) AuthService {
	return &authService{userRepo: userRepo, userService: userService}
}

func (s *authService) RegisterUser(user *models.User) (*models.User, error) {

	return s.userService.CreateUser(user, 1)
}

func (s *authService) LoginUser(identity, password string) (*models.User, string, string, error) {
	user := &models.User{}
	var err error

	if utils.ValidEmail(identity) {
		err = s.userRepo.FindField(user, EMAIL, identity)
		if err != nil {
			return nil, "", "", err
		}
	} else {
		err = s.userRepo.FindField(user, USERNAME, identity)
		if err != nil {
			return nil, "", "", err
		}
	}

	if user.ID == 0 {
		return nil, "", "", ErrUserNotFound
	}

	if !utils.ComparePassword(user.Password, password) {
		return nil, "", "", ErrInvalidUserPassword
	}

	accessToken, refreshToken, err := generateTokens(user)

	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *authService) RefreshAccessToken(refreshToken string) (string, string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod
		}
		return []byte(config.GetServerSettings().JWT.SecretKey), nil
	})

	if err != nil || !token.Valid {
		return "", "", ErrInvalidOrExpiredToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", "", ErrInvalidTokenClaims
	}

	username := claims[USERNAME].(string)

	if !ok {
		return "", "", ErrInvalidTokenClaims
	}

	var user = &models.User{}

	err = s.userRepo.FindField(user, USERNAME, username)
	if err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err := generateTokens(user)

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil

	//return newAccessToken, nil
}

func generateTokens(user *models.User) (string, string, error) {
	accessToken, err := generateJWT(user, accessTokenExpiryDuration)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := generateJWT(user, refreshTokenExpiryDuration)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func generateJWT(user *models.User, expiryDuration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["user_id"] = user.ID
	claims["email"] = user.Email
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(expiryDuration).Unix()

	t, err := token.SignedString([]byte(config.GetServerSettings().JWT.SecretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}
