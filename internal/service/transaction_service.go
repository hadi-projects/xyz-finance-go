package services

import (
	"errors"

	"github.com/hadi-projects/xyz-finance-go/internal/dto"
	"github.com/hadi-projects/xyz-finance-go/internal/entity"
	"github.com/hadi-projects/xyz-finance-go/internal/repository"
	"gorm.io/gorm"
)

type TransactionService interface {
	CreateTransaction(userId uint, req dto.CreateTransactionRequest) error
}

type transactionService struct {
	transactionRepo repository.TransactionRepository
	limitRepo       repository.LimitRepository
	db              *gorm.DB
}

func NewTransactionService(transactionRepo repository.TransactionRepository, limitRepo repository.LimitRepository, db *gorm.DB) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		limitRepo:       limitRepo,
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

		return nil
	})
}
