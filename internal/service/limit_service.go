package services

import (
	"github.com/hadi-projects/xyz-finance-go/internal/entity"
	"github.com/hadi-projects/xyz-finance-go/internal/repository"
)

type LimitService interface {
	GetConsumerLimit(userId uint) (*entity.TenorLimit, error)
}

type limitService struct {
	limitRepo repository.LimitRepository
}

func NewLimitService(limitRepo repository.LimitRepository) LimitService {
	return &limitService{
		limitRepo: limitRepo,
	}
}

func (s *limitService) GetConsumerLimit(userId uint) (*entity.TenorLimit, error) {
	limitData, err := s.limitRepo.FindByUserID(userId)
	if err != nil {
		return nil, err
	}
	return limitData, nil
}
