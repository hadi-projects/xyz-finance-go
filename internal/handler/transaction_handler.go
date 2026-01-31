package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/xyz-finance-go/internal/dto"
	services "github.com/hadi-projects/xyz-finance-go/internal/service"
	"github.com/hadi-projects/xyz-finance-go/pkg/logger"
)

type TransactionHandler struct {
	transactionService services.TransactionService
}

func NewTransactionHandler(transactionService services.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.GetUint("user_id")
	if err := h.transactionService.CreateTransaction(userId, req); err != nil {
		if err.Error() == "insufficient limit" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transaction created successfully"})

	logger.AuditLogger.Info().
		Str("action", "create_transaction").
		Uint("user_id", userId).
		Float64("otr", req.OTR).
		Float64("installment", req.InstallmentAmount).
		Int("tenor", req.Tenor).
		Msg("Transaction created")
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userId := c.GetUint("user_id")
	transactions, err := h.transactionService.GetTransactions(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions})
}
