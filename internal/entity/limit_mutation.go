package entity

import (
	"time"
)

type MutationAction string

const (
	MutationCreate MutationAction = "CREATE"
	MutationUpdate MutationAction = "UPDATE"
	MutationDelete MutationAction = "DELETE"
	MutationUsage  MutationAction = "USAGE"
)

type LimitMutation struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `json:"user_id"`
	TenorLimitID uint           `json:"tenor_limit_id"`
	OldAmount    float64        `json:"old_amount"`
	NewAmount    float64        `json:"new_amount"`
	Reason       string         `json:"reason"`
	Action       MutationAction `json:"action"` // CREATE, UPDATE, DELETE
	CreatedAt    time.Time      `json:"created_at"`
}

func (LimitMutation) TableName() string { return "limit_mutations" }
