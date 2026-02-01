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
	GetUserIDByLimitID(limitID uint) (uint, error)
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
	// Optimized: Use raw SQL JOIN instead of Association
	err := r.db.Raw(`
		SELECT tl.id, tl.tenor_month, tl.limit_amount, tl.created_at, tl.updated_at
		FROM tenor_limits tl
		INNER JOIN user_has_tenor_limit uhtl ON tl.id = uhtl.tenor_limit_id
		WHERE uhtl.user_id = ?
		ORDER BY tl.tenor_month ASC
	`, userId).Scan(&limits).Error
	if err != nil {
		return nil, err
	}
	return limits, nil
}

func (r *limitRepository) GetUserIDByLimitID(limitID uint) (uint, error) {
	var userID uint
	err := r.db.Raw("SELECT user_id FROM user_has_tenor_limit WHERE tenor_limit_id = ?", limitID).Scan(&userID).Error
	if err != nil {
		return 0, err
	}
	if userID == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return userID, nil
}

func (r *limitRepository) WithTx(tx *gorm.DB) LimitRepository {
	return &limitRepository{db: tx}
}
