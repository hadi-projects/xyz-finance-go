package entity

import "time"

// Role merepresentasikan peran user (e.g., Admin, User)
type Role struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"unique;type:varchar(50);not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relasi Many-to-Many: Satu Role punya banyak Permission
	Permissions []Permission `gorm:"many2many:role_has_permissions;" json:"permissions"`
}

func (Role) TableName() string {
	return "roles"
}
