package service

import (
	"errors"

	"com.mx/crud/internal/models"
	"com.mx/crud/internal/repository"
	"com.mx/crud/internal/utils"
	"github.com/gofiber/fiber/v2/log"
)

var (
	ErrInvalidUserPassword     = errors.New("invalid user/password")
	ErrInvalidEmailFormat      = errors.New("invalid email format")
	ErrUserNotFound            = errors.New("user not found")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidOrExpiredToken   = errors.New("invalid or expired refresh token")
	ErrInvalidTokenClaims      = errors.New("invalid token claims")

	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrFindingUserByEmail = errors.New("error finding user by email")
	ErrFindingUserByID    = errors.New("error finding user by id")
)

type UserService interface {
	GetAllUsers(page, pageSize int) ([]models.User, int64, error)
	CreateUser(user *models.User, auditUserID uint) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo}
}

func (s *userService) CreateUser(user *models.User, auditUserID uint) (*models.User, error) {

	if !utils.ValidEmail(user.Email) {
		log.Debug("Error parsing input: ")
		return nil, ErrInvalidEmailFormat
	}

	var userAux *models.User
	err := s.userRepo.FindField(userAux, "email", user.Email)

	if err != nil {
		log.Debug("Error finding user by email", err)
		return nil, ErrFindingUserByEmail
	}

	if userAux != nil {
		log.Debug("user already exists")
		return nil, ErrUserAlreadyExists
	}

	user.CreatedBy = auditUserID
	user.UpdatedBy = auditUserID

	if err := s.userRepo.Create(user); err != nil {
		log.Debug("Error creating user: ", err)
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
	var user *models.User
	err := s.userRepo.FindID(user, id)

	if user == nil {
		log.Debug("User not found")
		return nil, ErrUserNotFound
	}
	if err != nil {
		log.Debug("Error getting user: ", err)
		return nil, ErrFindingUserByID
	}

	return user, nil
}

func (s *userService) GetAllUsers(page, pageSize int) ([]models.User, int64, error) {
	return s.userRepo.FindAll(page, pageSize)
}

func (s *userService) UpdateUser(user *models.User) error {

	userAux, err := s.GetUserByID(user.ID)
	if err != nil {
		log.Debug("Error getting user:", err)
		return err
	}

	if user.Password != "" {
		hashedPassword, err := utils.GeneratePassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	} else {
		user.Password = userAux.Password
	}

	return s.userRepo.Update(user)
}

func (s *userService) DeleteUser(id uint) error {
	var user *models.User
	err := s.userRepo.FindID(user, id)
	if err != nil {
		return err
	}
	return s.userRepo.Delete(user)
}
