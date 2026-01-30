package entity

import "time"

// User adalah entity untuk Autentikasi (Login)
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"uniqueIndex;type:varchar(100);not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"` // JSON "-" agar password tidak ikut terkirim di API response
	RoleID    uint      `gorm:"not null" json:"role_id"`
	Role      Role      `gorm:"foreignKey:RoleID" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
