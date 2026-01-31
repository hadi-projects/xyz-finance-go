package services

import (
	"errors"

	"github.com/hadi-projects/xyz-finance-go/internal/dto"
	"github.com/hadi-projects/xyz-finance-go/internal/entity"
	"github.com/hadi-projects/xyz-finance-go/internal/repository"
	"github.com/hadi-projects/xyz-finance-go/pkg/logger"
	"gorm.io/gorm"
)

type TransactionService interface {
	CreateTransaction(userId uint, req dto.CreateTransactionRequest) error
}

type transactionService struct {
	transactionRepo repository.TransactionRepository
	limitRepo       repository.LimitRepository
	mutationRepo    repository.LimitMutationRepository
	db              *gorm.DB
}

func NewTransactionService(transactionRepo repository.TransactionRepository, limitRepo repository.LimitRepository, mutationRepo repository.LimitMutationRepository, db *gorm.DB) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		limitRepo:       limitRepo,
		mutationRepo:    mutationRepo,
		db:              db,
	}
}

func (s *transactionService) CreateTransaction(userId uint, req dto.CreateTransactionRequest) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Lock User Row (prevents race condition for this user)
		if err := tx.Exec("SELECT id FROM users WHERE id = ? FOR UPDATE", userId).Error; err != nil {
			return err
		}

		// 2. Use repositories with the transaction instance
		limitRepoTx := s.limitRepo.WithTx(tx)
		transactionRepoTx := s.transactionRepo.WithTx(tx)

		limits, err := limitRepoTx.FindByUserID(userId)
		if err != nil {
			return err
		}
		var limitAmount float64
		found := false
		for _, limit := range limits {
			if int(limit.TenorMonth) == req.Tenor {
				limitAmount = limit.LimitAmount
				found = true
				break
			}
		}

		if !found {
			return errors.New("limit not found for the requested tenor")
		}

		existingTransactions, err := transactionRepoTx.FindByUserID(userId)
		if err != nil {
			return err
		}

		var usedAmount float64
		for _, t := range existingTransactions {
			if t.Tenor == req.Tenor {
				usedAmount += t.OTR
			}
		}
		if usedAmount+req.OTR > limitAmount {
			return errors.New("insufficient limit")
		}

		transaction := &entity.Transaction{
			UserID:            userId,
			ContractNumber:    req.ContractNumber,
			OTR:               req.OTR,
			AdminFee:          req.AdminFee,
			InstallmentAmount: req.InstallmentAmount,
			InterestAmount:    req.InterestAmount,
			AssetName:         req.AssetName,
			Status:            "pending",
			Tenor:             req.Tenor,
		}

		if err := transactionRepoTx.Create(transaction); err != nil {
			return err
		}

		// Finds the Tenor Limit ID for the usage entry
		var limitID uint
		for _, limit := range limits {
			if int(limit.TenorMonth) == req.Tenor {
				limitID = uint(limit.ID)
				break
			}
		}

		// Log Usage Mutation
		mutation := &entity.LimitMutation{
			UserID:       userId,
			TenorLimitID: limitID,
			OldAmount:    limitAmount, // Current Limit Ceiling
			NewAmount:    limitAmount, // Current Limit Ceiling (Unchanged)
			Reason:       "Transaction Usage: " + req.ContractNumber,
			Action:       entity.MutationUsage,
		}

		if err := s.mutationRepo.WithTx(tx).Create(mutation); err != nil {
			return err
		}

		// Log to Audit File
		logger.AuditLogger.Info().
			Uint("user_id", userId).
			Uint("limit_id", limitID).
			Float64("amount", req.OTR).
			Str("contract_number", req.ContractNumber).
			Msg("Transaction Created (Limit Usage)")

		return nil
	})
}
