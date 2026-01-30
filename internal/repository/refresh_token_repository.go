package repository

import (
	"time"

	"github.com/hadi-projects/xyz-finance-go/internal/entity"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(token *entity.RefreshToken) error
	FindByToken(token string) (*entity.RefreshToken, error)
	RevokeByToken(token string) error
	RevokeAllByUserID(userID uint) error
	DeleteExpired() error
}

type refreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository creates a new refresh token repository instance
func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(token *entity.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *refreshTokenRepository) FindByToken(token string) (*entity.RefreshToken, error) {
	var refreshToken entity.RefreshToken
	err := r.db.Where("token = ? AND revoked = ? AND expires_at > ?",
		token, false, time.Now()).First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *refreshTokenRepository) RevokeByToken(token string) error {
	return r.db.Model(&entity.RefreshToken{}).
		Where("token = ?", token).
		Update("revoked", true).Error
}

func (r *refreshTokenRepository) RevokeAllByUserID(userID uint) error {
	return r.db.Model(&entity.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("revoked", true).Error
}

func (r *refreshTokenRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).
		Delete(&entity.RefreshToken{}).Error
}
