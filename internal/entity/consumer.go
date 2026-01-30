package entity

import (
	"time"
)

type Consumer struct {
	ID           uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	NIK          string  `gorm:"uniqueIndex;type:varchar(16);not null" json:"nik"`
	FullName     string  `gorm:"type:varchar(100);not null" json:"full_name"`
	LegalName    string  `gorm:"type:varchar(100);not null" json:"legal_name"`
	PlaceOfBirth string  `gorm:"type:varchar(50)" json:"place_of_birth"`
	DateOfBirth  string  `gorm:"type:date" json:"date_of_birth"` // Format YYYY-MM-DD
	Salary       float64 `gorm:"type:decimal(15,2)" json:"salary"`
	KTPImage     string  `gorm:"type:varchar(255)" json:"ktp_image"`
	SelfieImage  string  `gorm:"type:varchar(255)" json:"selfie_image"`

	Limits []TenorLimit `gorm:"foreignKey:ConsumerID" json:"limits"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Consumer) TableName() string { return "consumers" }
