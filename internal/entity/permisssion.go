package entity

import "time"

// Permission merepresentasikan hak akses granular
type Permission struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"unique;type:varchar(100);not null" json:"name"` // e.g. "transaction.create"

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Permission) TableName() string {
	return "permissions"
}
