package entity

import (
	"time"
)

type Tenor int

const (
	Tenor1 Tenor = 1
	Tenor2 Tenor = 2
	Tenor3 Tenor = 3
	Tenor6 Tenor = 6
)

type TenorLimit struct {
	ID          uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	TenorMonth  Tenor   `gorm:"not null;comment:'1, 2, 3, or 6'" json:"tenor_month"`
	LimitAmount float64 `gorm:"type:decimal(15,2);default:0" json:"limit_amount"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TenorLimit) TableName() string { return "tenor_limits" }
