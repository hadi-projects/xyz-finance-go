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
		limits, err := s.limitRepo.FindByUserID(userId)
		if err != nil {
			return err
		}
		var limitAmount float64
		found := false
		for _, limit := range limits {
			if limit.TenorMonth == req.Tenor {
				limitAmount = limit.LimitAmount
				found = true
				break
			}
		}

		if !found {
			return errors.New("limit not found for the requested tenor")
		}

		existingTransactions, err := s.transactionRepo.FindByUserID(userId)
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

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		return nil
	})
}
