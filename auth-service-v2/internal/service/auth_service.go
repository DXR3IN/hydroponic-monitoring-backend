package service

import (
	"errors"

	models "github.com/DXR3IN/auth-service-v2/internal/domain"
	"github.com/DXR3IN/auth-service-v2/internal/repository"
	"github.com/DXR3IN/auth-service-v2/internal/utils"
	"github.com/DXR3IN/auth-service-v2/pkg/logger"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidSubjectID   = errors.New("token subject is not a valid user ID")
)

type AuthService struct {
	repo repository.UserRepository
	jwt  *utils.JWTManager
}

func NewAuthService(r repository.UserRepository, jwt *utils.JWTManager) *AuthService {
	return &AuthService{repo: r, jwt: jwt}
}

func (s *AuthService) Register(name, email, password string) (string, error) {
	// check if email already used
	ex, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if ex != nil {
		return "", ErrUserExists
	}

	hashed, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	u := &repository.User{Name: name, Email: email, Password: hashed}
	if err := s.repo.Create(u); err != nil {
		return "", err
	}

	token, err := s.jwt.Generate(u.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	u, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", err
	}
	if u == nil {
		return "", ErrInvalidCredentials
	}

	if err := utils.CheckPasswordHash(password, u.Password); err != nil {
		return "", ErrInvalidCredentials
	}

	// Generate token
	token, err := s.jwt.Generate(u.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) GetUserDataByID(userID string) (*models.User, error) {
	u, err := s.repo.FindByID(userID)
	if err != nil {
		logger.ErrorLogger(err)
		return nil, err
	}
	return u, nil
}

func (s *AuthService) VerifyToken(token string) (string, error) {
	claims, err := s.jwt.Verify(token)
	if err != nil {
		return "", err
	}

	return claims.Subject, nil
}

func (s *AuthService) UpdateName(userID, newName string) error {
	return s.repo.EditNameByID(userID, newName)
}

func (s *AuthService) UpdatePassword(userID, newPassword string) error {
	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	return s.repo.EditPasswordByID(userID, hashed)
}
