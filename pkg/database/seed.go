package database

import (
	"errors"
	"fmt"

	"github.com/hadi-projects/xyz-finance-go/internal/entity"
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
			fmt.Printf("⚠️ Failed to create role %s: %v\n", roleName, err)
		}
	} else if err == nil {
		// Jika Role sudah ada, update permission-nya (sync)
		// db.Model(&role).Association("Permissions").Replace(perms)
		// (Opsional: di-uncomment jika ingin permission selalu reset setiap restart)
	}
}
