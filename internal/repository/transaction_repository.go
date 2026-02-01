package repository

import (
	"github.com/hadi-projects/xyz-finance-go/internal/entity"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(transaction *entity.Transaction) error
	FindByUserID(userId uint) ([]entity.Transaction, error)
	FindByUserIDPaginated(userId uint, offset, limit int) ([]entity.Transaction, int64, error)
	FindAll() ([]entity.Transaction, error)
	FindAllPaginated(offset, limit int) ([]entity.Transaction, int64, error)
	WithTx(tx *gorm.DB) TransactionRepository
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(transaction *entity.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) FindByUserID(userId uint) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := r.db.Where("user_id = ?", userId).Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) FindByUserIDPaginated(userId uint, offset, limit int) ([]entity.Transaction, int64, error) {
	var transactions []entity.Transaction
	var total int64

	if err := r.db.Model(&entity.Transaction{}).Where("user_id = ?", userId).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Select("id", "user_id", "contract_number", "otr", "admin_fee", "installment_amount", "interest_amount", "asset_name", "status", "tenor", "created_at", "updated_at").
		Where("user_id = ?", userId).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&transactions).Error

	return transactions, total, err
}

func (r *transactionRepository) FindAll() ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	err := r.db.Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) FindAllPaginated(offset, limit int) ([]entity.Transaction, int64, error) {
	var transactions []entity.Transaction
	var total int64

	// Count total
	if err := r.db.Model(&entity.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated data with Select for specific columns
	err := r.db.Select("id", "user_id", "contract_number", "otr", "admin_fee", "installment_amount", "interest_amount", "asset_name", "status", "tenor", "created_at", "updated_at").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&transactions).Error

	return transactions, total, err
}

func (r *transactionRepository) WithTx(tx *gorm.DB) TransactionRepository {
	return &transactionRepository{db: tx}
}
