package services

import (
	"errors"

	"github.com/hadi-projects/xyz-finance-go/internal/dto"
	"github.com/hadi-projects/xyz-finance-go/internal/entity"
	"github.com/hadi-projects/xyz-finance-go/internal/repository"
)

type LimitService interface {
	GetLimitsByUserID(userId uint) ([]entity.TenorLimit, error)
	CreateLimit(req dto.CreateLimitRequest) error
	UpdateLimit(id uint, req dto.UpdateLimitRequest) error
	DeleteLimit(id uint) error
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
	// Validate Tenor
	validTenors := map[int]bool{1: true, 2: true, 3: true, 6: true}
	if !validTenors[req.TenorMonth] {
		return errors.New("invalid tenor month: must be 1, 2, 3, or 6")
	}

	// Validate User Exists
	if _, err := s.userRepo.FindByID(req.TargetUserID); err != nil {
		return errors.New("target user not found")
	}

	limit := &entity.TenorLimit{
		TenorMonth:  entity.Tenor(req.TenorMonth),
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

func (s *limitService) UpdateLimit(id uint, req dto.UpdateLimitRequest) error {
	validTenors := map[int]bool{1: true, 2: true, 3: true, 6: true}
	if !validTenors[req.TenorMonth] {
		return errors.New("invalid tenor month: must be 1, 2, 3, or 6")
	}

	limit, err := s.limitRepo.FindByID(id)
	if err != nil {
		return errors.New("limit not found")
	}

	limit.TenorMonth = entity.Tenor(req.TenorMonth)
	limit.LimitAmount = req.LimitAmount

	return s.limitRepo.Update(limit)
}

func (s *limitService) DeleteLimit(id uint) error {
	if _, err := s.limitRepo.FindByID(id); err != nil {
		return errors.New("limit not found")
	}
	return s.limitRepo.Delete(id)
}
