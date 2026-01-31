package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hadi-projects/xyz-finance-go/internal/dto"
	"github.com/hadi-projects/xyz-finance-go/internal/handler"
	"github.com/hadi-projects/xyz-finance-go/internal/service/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTransactionHandler_CreateTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxService := mock.NewMockTransactionService(ctrl)
	txHandler := handler.NewTransactionHandler(mockTxService)

	t.Run("Success", func(t *testing.T) {
		req := dto.CreateTransactionRequest{
			ContractNumber:    "CTR-001",
			OTR:               10000,
			AdminFee:          500,
			InstallmentAmount: 1100,
			InterestAmount:    100,
			AssetName:         "Item1",
			Tenor:             1,
		}
		body, _ := json.Marshal(req)
		userId := uint(1)

		mockTxService.EXPECT().CreateTransaction(userId, req).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/transaction/", bytes.NewBuffer(body))
		c.Set("user_id", userId)

		txHandler.CreateTransaction(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		req := dto.CreateTransactionRequest{
			ContractNumber:    "CTR-001",
			OTR:               10000,
			AdminFee:          500,
			InstallmentAmount: 1100,
			InterestAmount:    100,
			AssetName:         "Item1",
			Tenor:             1,
		}
		body, _ := json.Marshal(req)
		userId := uint(1)

		mockTxService.EXPECT().CreateTransaction(userId, req).Return(errors.New("insufficient limit"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/transaction/", bytes.NewBuffer(body))
		c.Set("user_id", userId)

		txHandler.CreateTransaction(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		req := dto.CreateTransactionRequest{}
		body, _ := json.Marshal(req)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/transaction/", bytes.NewBuffer(body))
		// No user_id set

		txHandler.CreateTransaction(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
