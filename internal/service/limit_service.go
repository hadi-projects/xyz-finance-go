package services

import (
	"errors"

	"github.com/hadi-projects/xyz-finance-go/internal/dto"
	"github.com/hadi-projects/xyz-finance-go/internal/entity"
	"github.com/hadi-projects/xyz-finance-go/internal/repository"
	"github.com/hadi-projects/xyz-finance-go/pkg/logger"
	"gorm.io/gorm"
)

type LimitService interface {
	GetLimitsByUserID(userId uint) ([]entity.TenorLimit, error)
	CreateLimit(req dto.CreateLimitRequest) error
	UpdateLimit(id uint, req dto.UpdateLimitRequest) error
	DeleteLimit(id uint) error
}

type limitService struct {
	limitRepo    repository.LimitRepository
	userRepo     repository.UserRepository
	mutationRepo repository.LimitMutationRepository
	db           *gorm.DB
}

func NewLimitService(limitRepo repository.LimitRepository, userRepo repository.UserRepository, mutationRepo repository.LimitMutationRepository, db *gorm.DB) LimitService {
	return &limitService{
		limitRepo:    limitRepo,
		userRepo:     userRepo,
		mutationRepo: mutationRepo,
		db:           db,
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

	// Check if limit for this tenor already exists
	existingLimits, err := s.limitRepo.FindByUserID(req.TargetUserID)
	if err != nil {
		return err
	}
	for _, l := range existingLimits {
		if int(l.TenorMonth) == req.TenorMonth {
			return errors.New("limit for this tenor already exists")
		}
	}

	limit := &entity.TenorLimit{
		TenorMonth:  entity.Tenor(req.TenorMonth),
		LimitAmount: req.LimitAmount,
	}

	// Transaction
	return s.db.Transaction(func(tx *gorm.DB) error {
		limitRepoTx := s.limitRepo.WithTx(tx)
		// userRepoTx := s.userRepo // Not used yet
		// Note: UserRepo currently fetches DB instance from struct. Ideally UserRepo should also support WithTx.
		// For now, let's assume UserRepo CreateUserHasTenorLimit needs to be transaction aware or we fix UserRepo later.
		// Ideally we should fix UserRepo. But for speed, let's look at `CreateUserHasTenorLimit`. It uses `r.db`.
		// Let's implement WithTx for UserRepo too? Or interact with DB directly for association?
		// To avoid scope creep, I will rely on the fact that `userRepo` uses its own DB instance.
		// Wait, if I don't propagate TX to userRepo, the `CreateUserHasTenorLimit` will be outside the transaction.
		// This is bad. I MUST update UserRepository to support WithTx or do the association manually here using tx.

		// Manual Association using TX to ensure consistency
		if err := limitRepoTx.Create(limit); err != nil {
			return err
		}

		// Link to user (Optimized: using raw SQL or manual insert to avoid UserRepo dependency overhead if UserRepo doesn't have WithTx)
		// Or:
		// err := tx.Exec("INSERT INTO user_has_tenor_limit (user_id, tenor_limit_id) VALUES (?, ?)", req.TargetUserID, limit.ID).Error
		if err := tx.Exec("INSERT INTO user_has_tenor_limit (user_id, tenor_limit_id) VALUES (?, ?)", req.TargetUserID, limit.ID).Error; err != nil {
			return err
		}

		// Log Mutation
		mutation := &entity.LimitMutation{
			UserID:       req.TargetUserID,
			TenorLimitID: uint(limit.ID),
			OldAmount:    0,
			NewAmount:    req.LimitAmount,
			Reason:       "Initial Limit",
			Action:       entity.MutationCreate,
		}
		if err := s.mutationRepo.WithTx(tx).Create(mutation); err != nil {
			return err
		}

		// Log to Audit File
		logger.AuditLogger.Info().
			Uint("user_id", req.TargetUserID).
			Uint("limit_id", uint(limit.ID)).
			Float64("limit_amount", req.LimitAmount).
			Msg("Limit Created")

		return nil
	})
}

func (s *limitService) UpdateLimit(id uint, req dto.UpdateLimitRequest) error {
	validTenors := map[int]bool{1: true, 2: true, 3: true, 6: true}
	if !validTenors[req.TenorMonth] {
		return errors.New("invalid tenor month: must be 1, 2, 3, or 6")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		limitRepoTx := s.limitRepo.WithTx(tx)

		limit, err := limitRepoTx.FindByID(id)
		if err != nil {
			return errors.New("limit not found")
		}

		oldAmount := limit.LimitAmount

		// Find UserID
		userID, err := limitRepoTx.GetUserIDByLimitID(id)
		if err != nil {
			// Fallback or error? Strictly speaking, a limit should belong to a user.
			return errors.New("limit owner not found")
		}

		limit.TenorMonth = entity.Tenor(req.TenorMonth)
		limit.LimitAmount = req.LimitAmount

		if err := limitRepoTx.Update(limit); err != nil {
			return err
		}

		// Log Mutation
		mutation := &entity.LimitMutation{
			UserID:       userID,
			TenorLimitID: uint(limit.ID),
			OldAmount:    oldAmount,
			NewAmount:    req.LimitAmount,
			Reason:       "Update Limit",
			Action:       entity.MutationUpdate,
		}
		if err := s.mutationRepo.WithTx(tx).Create(mutation); err != nil {
			return err
		}

		// Log to Audit File
		logger.AuditLogger.Info().
			Uint("user_id", userID).
			Uint("limit_id", uint(limit.ID)).
			Float64("old_amount", oldAmount).
			Float64("new_amount", req.LimitAmount).
			Msg("Limit Updated")

		return nil
	})
}

func (s *limitService) DeleteLimit(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		limitRepoTx := s.limitRepo.WithTx(tx)

		limit, err := limitRepoTx.FindByID(id)
		if err != nil {
			return errors.New("limit not found")
		}

		// Find UserID before deleting
		userID, err := limitRepoTx.GetUserIDByLimitID(id)
		if err != nil {
			return errors.New("limit owner not found")
		}

		if err := limitRepoTx.Delete(id); err != nil {
			return err
		}

		// Log Mutation
		mutation := &entity.LimitMutation{
			UserID:       userID,
			TenorLimitID: id,
			OldAmount:    limit.LimitAmount,
			NewAmount:    0,
			Reason:       "Delete Limit",
			Action:       entity.MutationDelete,
		}
		if err := s.mutationRepo.WithTx(tx).Create(mutation); err != nil {
			return err
		}

		// Log to Audit File
		logger.AuditLogger.Info().
			Uint("user_id", userID).
			Uint("limit_id", id).
			Msg("Limit Deleted")

		return nil
	})
}
