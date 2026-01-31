package entity

import "time"

type Transaction struct {
	ID                uint64  `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID            uint    `gorm:"not null" json:"user_id"`
	User              User    `gorm:"foreignKey:UserID" json:"-"`
	ContractNumber    string  `gorm:"uniqueIndex;type:varchar(50);not null" json:"contract_number"`
	OTR               float64 `gorm:"type:decimal(15,2);not null" json:"otr"`
	AdminFee          float64 `gorm:"type:decimal(15,2);not null" json:"admin_fee"`
	InstallmentAmount float64 `gorm:"type:decimal(15,2);not null" json:"installment_amount"`
	InterestAmount    float64 `gorm:"type:decimal(15,2);not null" json:"interest_amount"`
	AssetName         string  `gorm:"type:varchar(255);not null" json:"asset_name"`
	Status            string  `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, approved, rejected
	Tenor             int     `gorm:"type:int;not null" json:"tenor"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Transaction) TableName() string {
	return "transactions"
}
