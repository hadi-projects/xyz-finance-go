package services

import (
	"errors"
	"fmt"

	"github.com/hadi-projects/xyz-finance-go/internal/entity"
	"github.com/hadi-projects/xyz-finance-go/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(email, password string) (*entity.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(
	userRepo repository.UserRepository,
) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) Register(email, password string) (*entity.User, error) {
	_, err := s.userRepo.FindByEmail(email)
	if err == nil {
		return nil, fmt.Errorf("user with this email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &entity.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}
