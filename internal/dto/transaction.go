package dto

type CreateTransactionRequest struct {
	ContractNumber    string  `json:"contract_number" binding:"required"`
	OTR               float64 `json:"otr" binding:"required"`
	AdminFee          float64 `json:"admin_fee" binding:"required"`
	InstallmentAmount float64 `json:"installment_amount" binding:"required"`
	InterestAmount    float64 `json:"interest_amount" binding:"required"`
	AssetName         string  `json:"asset_name" binding:"required"`
	Tenor             int     `json:"tenor" binding:"required"`
}
