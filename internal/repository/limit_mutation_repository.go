package repository

import (
	"github.com/hadi-projects/xyz-finance-go/internal/entity"
	"gorm.io/gorm"
)

type LimitMutationRepository interface {
	Create(mutation *entity.LimitMutation) error
	WithTx(tx *gorm.DB) LimitMutationRepository
}

type limitMutationRepository struct {
	db *gorm.DB
}

func NewLimitMutationRepository(db *gorm.DB) LimitMutationRepository {
	return &limitMutationRepository{db: db}
}

func (r *limitMutationRepository) Create(mutation *entity.LimitMutation) error {
	return r.db.Create(mutation).Error
}

func (r *limitMutationRepository) WithTx(tx *gorm.DB) LimitMutationRepository {
	return &limitMutationRepository{db: tx}
}
