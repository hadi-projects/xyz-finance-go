package handler

import (
	services "github.com/hadi-projects/xyz-finance-go/internal/service"
)

type LimitHandler struct {
	limitService services.LimitService
}

// NewLimitHandler creates a new limit handler instance
func NewLimitHandler(limitService services.LimitService) *LimitHandler {
	return &LimitHandler{
		limitService: limitService,
	}
}
