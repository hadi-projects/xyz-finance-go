package repository

import (
	"github.com/hadi-projects/xyz-finance-go/internal/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entity.User) error
	FindByID(id uint) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id uint) error
	CreateUserHasTenorLimit(userId uint, limitID uint) error
	GetLimitsByUserID(userID uint) ([]entity.TenorLimit, error)
	FindAllWithLimits() ([]entity.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) error {
	if user.RoleID == 0 {
		user.RoleID = 1
	}
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Role.Permissions").Preload("Consumer").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&entity.User{}, id).Error
}

func (r *userRepository) CreateUserHasTenorLimit(userId uint, limitID uint) error {
	var user entity.User
	if err := r.db.First(&user, userId).Error; err != nil {
		return err
	}
	var limit entity.TenorLimit
	if err := r.db.First(&limit, limitID).Error; err != nil {
		return err
	}

	if err := r.db.Model(&user).Association("TenorLimit").Append(&limit); err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetLimitsByUserID(userID uint) ([]entity.TenorLimit, error) {
	var limits []entity.TenorLimit

	err := r.db.Model(&entity.User{ID: userID}).Association("TenorLimit").Find(&limits)

	return limits, err
}

func (r *userRepository) FindAllWithLimits() ([]entity.User, error) {
	var users []entity.User
	// Preload TenorLimit to get the limits for each user
	err := r.db.Preload("TenorLimit").Find(&users).Error
	return users, err
}
