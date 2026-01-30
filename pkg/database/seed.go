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
	seedRole(db, "admin", []entity.Permission{})
	seedRole(db, "user", []entity.Permission{})

	fmt.Println("✅ RBAC Seeding Completed!")
}

func seedRole(db *gorm.DB, roleName string, perms []entity.Permission) {
	var role entity.Role
	err := db.Where("name = ?", roleName).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		role = entity.Role{Name: roleName, Permissions: perms}
		if err := db.Create(&role).Error; err != nil {
			fmt.Printf("Failed to create role %s: %v\n", roleName, err)
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
	user := entity.User{Email: email, Password: string(password), RoleID: roleId}
	if err := repository.NewUserRepository(db).Create(&user); err != nil {
		fmt.Printf("Failed to create user %s: %v\n", email, err)
	}
}

// func SeedConsumerLimit(db *gorm.DB) {
// 	// budi
// 	seedLimit(db, 2, 1, 100000)
// 	seedLimit(db, 2, 2, 200000)
// 	seedLimit(db, 2, 3, 500000)
// 	seedLimit(db, 2, 6, 700000)

// 	// annisa
// 	seedLimit(db, 3, 1, 1000000)
// 	seedLimit(db, 3, 2, 1200000)
// 	seedLimit(db, 3, 3, 1500000)
// 	seedLimit(db, 3, 6, 2000000)

// 	fmt.Println("Consumer Limit Seeding Completed!")
// }

// func seedLimit(db *gorm.DB, tenor uint64, limitAmount uint64) {
// 	limit := entity.TenorLimit{TenorMonth: int(tenor), LimitAmount: float64(limitAmount)}
// 	if err := repository.NewLimitRepository(db).Create(&limit); err != nil {
// 		fmt.Printf("Failed to create limit %d: %v\n", tenor, err)
// 	}
// }
