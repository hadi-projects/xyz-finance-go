package repository

import (
	"github.com/hadi-projects/xyz-finance-go/internal/entity"
	"gorm.io/gorm"
)

type LimitRepository interface {
	Create(user *entity.TenorLimit) error
	FindByID(id uint) (*entity.TenorLimit, error)
	FindByEmail(email string) (*entity.TenorLimit, error)
	Update(user *entity.TenorLimit) error
	Delete(id uint) error
	FindByUserID(userId uint) ([]entity.TenorLimit, error)
	WithTx(tx *gorm.DB) LimitRepository
}

type limitRepository struct {
	db *gorm.DB
}

func NewLimitRepository(db *gorm.DB) LimitRepository {
	return &limitRepository{db: db}
}

func (r *limitRepository) Create(limit *entity.TenorLimit) error {
	return r.db.Create(limit).Error
}

func (r *limitRepository) FindByID(id uint) (*entity.TenorLimit, error) {
	var limit entity.TenorLimit
	err := r.db.First(&limit, id).Error
	if err != nil {
		return nil, err
	}
	return &limit, nil
}

func (r *limitRepository) FindByEmail(email string) (*entity.TenorLimit, error) {
	var limit entity.TenorLimit
	err := r.db.Where("email = ?", email).First(&limit).Error
	if err != nil {
		return nil, err
	}
	return &limit, nil
}

func (r *limitRepository) Update(limit *entity.TenorLimit) error {
	return r.db.Save(limit).Error
}

func (r *limitRepository) Delete(id uint) error {
	return r.db.Delete(&entity.TenorLimit{}, id).Error
}

func (r *limitRepository) FindByUserID(userId uint) ([]entity.TenorLimit, error) {
	var limits []entity.TenorLimit
	err := r.db.Model(&entity.User{ID: userId}).Association("TenorLimit").Find(&limits)
	if err != nil {
		return nil, err
	}
	return limits, nil
}

func (r *limitRepository) WithTx(tx *gorm.DB) LimitRepository {
	return &limitRepository{db: tx}
}
