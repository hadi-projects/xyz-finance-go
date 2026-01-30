package services

import (
	"github.com/hadi-projects/xyz-finance-go/internal/repository"
)

type LimitService interface {
	GetConsumerLimit(limit string) (string, error)
}

type limitService struct {
	limitRepo repository.LimitRepository
}

func NewLimitService(limitRepo repository.LimitRepository) LimitService {
	return &limitService{
		limitRepo: limitRepo,
	}
}

func (s *limitService) GetConsumerLimit(limit string) (string, error) {
	return limit, nil
}
