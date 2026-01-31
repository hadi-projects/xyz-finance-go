package database

import (
	"errors"
	"fmt"

	"github.com/hadi-projects/xyz-finance-go/internal/entity"
	"github.com/hadi-projects/xyz-finance-go/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedRBAC(db *gorm.DB) {
	seedRole(db, "admin", []entity.Permission{{Name: "create-limit"}, {Name: "delete-limit"}})
	seedRole(db, "user", []entity.Permission{{Name: "get-limit"}, {Name: "create-transaction"}})

	fmt.Println("✅ RBAC Seeding Completed!")
}

func seedRole(db *gorm.DB, roleName string, perms []entity.Permission) {
	var finalPerms []entity.Permission
	for _, p := range perms {
		var perm entity.Permission
		if err := db.Where("name = ?", p.Name).FirstOrCreate(&perm, entity.Permission{Name: p.Name}).Error; err != nil {
			fmt.Printf("Failed to seed permission %s: %v\n", p.Name, err)
			continue
		}
		finalPerms = append(finalPerms, perm)
	}

	var role entity.Role
	err := db.Where("name = ?", roleName).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		role = entity.Role{Name: roleName, Permissions: finalPerms}
		if err := db.Create(&role).Error; err != nil {
			fmt.Printf("Failed to create role %s: %v\n", roleName, err)
		}
	} else {
		if err := db.Model(&role).Association("Permissions").Replace(finalPerms); err != nil {
			fmt.Printf("Failed to update permissions for role %s: %v\n", roleName, err)
		}
	}
}

func SeedUser(db *gorm.DB) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("pAsswj@123"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("failed to hash password: %v\n", err)
	}

	hashedPassword2, err := bcrypt.GenerateFromPassword([]byte("pAsswj@1873"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("failed to hash password: %v\n", err)
	}

	hashedPassword3, err := bcrypt.GenerateFromPassword([]byte("pAsswj@1763"), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("failed to hash password: %v\n", err)
	}

	seedUser(db, "admin@mail.com", hashedPassword, 1)
	seedUser(db, "budi@mail.com", hashedPassword2, 2)
	seedUser(db, "annisa@mail.com", hashedPassword3, 2)

	fmt.Println("User Seeding Completed!")
}

func seedUser(db *gorm.DB, email string, password []byte, roleId uint) {
	var user entity.User
	if err := db.Where("email = ?", email).First(&user).Error; err == nil {
		return // User already exists
	}

	user = entity.User{Email: email, Password: string(password), RoleID: roleId}
	if err := repository.NewUserRepository(db).Create(&user); err != nil {
		fmt.Printf("Failed to create user %s: %v\n", email, err)
	}
}

func SeedConsumerLimit(db *gorm.DB) {
	// budi
	seedLimit(db, 2, 1, 100000)
	seedLimit(db, 2, 2, 200000)
	seedLimit(db, 2, 3, 500000)
	seedLimit(db, 2, 6, 700000)

	// annisa
	seedLimit(db, 3, 1, 1000000)
	seedLimit(db, 3, 2, 1200000)
	seedLimit(db, 3, 3, 1500000)
	seedLimit(db, 3, 6, 2000000)

	fmt.Println("Consumer Limit Seeding Completed!")
}

func seedLimit(db *gorm.DB, userId uint, tenor int, limitAmount float64) {
	// Check if user already has this limit
	var count int64
	db.Table("tenor_limits").
		Joins("JOIN user_has_tenor_limit ON user_has_tenor_limit.tenor_limit_id = tenor_limits.id").
		Where("user_has_tenor_limit.user_id = ? AND tenor_limits.tenor_month = ?", userId, tenor).
		Count(&count)

	if count > 0 {
		return // Limit already exists
	}

	limit := entity.TenorLimit{TenorMonth: entity.Tenor(tenor), LimitAmount: limitAmount}
	if err := repository.NewLimitRepository(db).Create(&limit); err != nil {
		fmt.Printf("Failed to create limit %d: %v\n", tenor, err)
	}

	// create user has tenor limit
	if err := repository.NewUserRepository(db).CreateUserHasTenorLimit(userId, uint(limit.ID)); err != nil {
		fmt.Printf("Failed to update user has tenor limit %d: %v\n", tenor, err)
	}
}

func SeedConsumer(db *gorm.DB) {
	seedConsumerData(db, 2, "1234567890123456", "Budi Santoso", "Budi Santoso", "Jakarta", "1990-01-01", 10000000)
	seedConsumerData(db, 3, "6543210987654321", "Annisa Putri", "Annisa Putri", "Bandung", "1992-05-15", 15000000)

	fmt.Println("Consumer Seeding Completed!")
}

func seedConsumerData(db *gorm.DB, userId uint, nik, fullName, legalName, pob, dob string, salary float64) {
	consumer := entity.Consumer{
		UserID:       userId,
		NIK:          nik,
		FullName:     fullName,
		LegalName:    legalName,
		PlaceOfBirth: pob,
		DateOfBirth:  dob,
		Salary:       salary,
		KTPImage:     "ktp_placeholder.jpg",
		SelfieImage:  "selfie_placeholder.jpg",
	}

	var existing entity.Consumer
	err := db.Where("user_id = ?", userId).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := db.Create(&consumer).Error; err != nil {
			fmt.Printf("Failed to create consumer for user %d: %v\n", userId, err)
		}
	} else {
		consumer.ID = existing.ID
		consumer.CreatedAt = existing.CreatedAt
		if err := db.Save(&consumer).Error; err != nil {
			fmt.Printf("Failed to update consumer for user %d: %v\n", userId, err)
		}
	}
}
