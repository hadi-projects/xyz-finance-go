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
	"github.com/hadi-projects/xyz-finance-go/internal/entity"
	"github.com/hadi-projects/xyz-finance-go/internal/handler"
	"github.com/hadi-projects/xyz-finance-go/internal/service/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLimitHandler_CreateLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLimitService := mock.NewMockLimitService(ctrl)
	limitHandler := handler.NewLimitHandler(mockLimitService)

	t.Run("Success", func(t *testing.T) {
		req := dto.CreateLimitRequest{
			TargetUserID: 1,
			TenorMonth:   12,
			LimitAmount:  1000000,
		}
		body, _ := json.Marshal(req)

		mockLimitService.EXPECT().CreateLimit(req).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/limit/", bytes.NewBuffer(body))

		limitHandler.CreateLimit(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		req := dto.CreateLimitRequest{TargetUserID: 1, TenorMonth: 12, LimitAmount: 100}
		body, _ := json.Marshal(req)

		mockLimitService.EXPECT().CreateLimit(req).Return(errors.New("service error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/limit/", bytes.NewBuffer(body))

		limitHandler.CreateLimit(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestLimitHandler_GetLimits(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLimitService := mock.NewMockLimitService(ctrl)
	limitHandler := handler.NewLimitHandler(mockLimitService)

	t.Run("Success", func(t *testing.T) {
		userId := uint(1)
		limits := []entity.TenorLimit{{ID: 1, TenorMonth: 12, LimitAmount: 10000}}

		mockLimitService.EXPECT().GetLimitsByUserID(userId).Return(limits, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/limit/", nil)
		c.Set("user_id", userId) // Simulate auth middleware

		limitHandler.GetLimits(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestLimitHandler_DeleteLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLimitService := mock.NewMockLimitService(ctrl)
	limitHandler := handler.NewLimitHandler(mockLimitService)

	t.Run("Success", func(t *testing.T) {
		limitID := 123
		mockLimitService.EXPECT().DeleteLimit(uint(limitID)).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "123"}}
		c.Request, _ = http.NewRequest("DELETE", "/api/limit/123", nil)

		limitHandler.DeleteLimit(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
