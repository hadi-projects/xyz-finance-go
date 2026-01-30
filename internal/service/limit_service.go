package services

import (
	"github.com/hadi-projects/xyz-finance-go/internal/dto"
	"github.com/hadi-projects/xyz-finance-go/internal/entity"
	"github.com/hadi-projects/xyz-finance-go/internal/repository"
)

type LimitService interface {
	GetLimitsByUserID(userId uint) ([]entity.TenorLimit, error)
	CreateLimit(req dto.CreateLimitRequest) error
}

type limitService struct {
	limitRepo repository.LimitRepository
	userRepo  repository.UserRepository
}

func NewLimitService(limitRepo repository.LimitRepository, userRepo repository.UserRepository) LimitService {
	return &limitService{
		limitRepo: limitRepo,
		userRepo:  userRepo,
	}
}

func (s *limitService) GetLimitsByUserID(userId uint) ([]entity.TenorLimit, error) {
	limits, err := s.limitRepo.FindByUserID(userId)
	if err != nil {
		return nil, err
	}
	return limits, nil
}

func (s *limitService) CreateLimit(req dto.CreateLimitRequest) error {
	limit := &entity.TenorLimit{
		TenorMonth:  req.TenorMonth,
		LimitAmount: req.LimitAmount,
	}

	if err := s.limitRepo.Create(limit); err != nil {
		return err
	}

	// Link to user
	// Note: limit.ID is populated after Create
	if err := s.userRepo.CreateUserHasTenorLimit(req.TargetUserID, uint(limit.ID)); err != nil {
		return err
	}

	return nil
}
