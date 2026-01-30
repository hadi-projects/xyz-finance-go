package dto

type CreateLimitRequest struct {
	TargetUserID uint    `json:"target_user_id" binding:"required"`
	TenorMonth   int     `json:"tenor_month" binding:"required"`
	LimitAmount  float64 `json:"limit_amount" binding:"required"`
}
