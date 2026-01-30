package entity

import (
	"time"
)

type TenorLimit struct {
	ID          uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	TenorMonth  int     `gorm:"not null;comment:'1, 2, 3, or 6'" json:"tenor_month"`
	LimitAmount float64 `gorm:"type:decimal(15,2);default:0" json:"limit_amount"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TenorLimit) TableName() string { return "tenor_limits" }
