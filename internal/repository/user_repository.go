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
	CreateUserHasTenorLimit(userId uint, tenor int) error
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
	err := r.db.Preload("Profile").First(&user, id).Error
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

func (r *userRepository) CreateUserHasTenorLimit(userId uint, tenor int) error {
	var user entity.User
	if err := r.db.First(&user, userId).Error; err != nil {
		panic(err)
	}
	var limit entity.TenorLimit
	if err := r.db.Where("tenor_month = ?", tenor).First(&limit).Error; err != nil {
		panic(err)
	}

	err := r.db.Model(&user).Association("TenorLimit").Append(&limit)
	if err != nil {
		panic(err)
	}

	return nil
}
