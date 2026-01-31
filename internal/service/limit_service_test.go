package services_test

import (
	"errors"
	"testing"

	"github.com/hadi-projects/xyz-finance-go/internal/dto"
	"github.com/hadi-projects/xyz-finance-go/internal/entity"
	"github.com/hadi-projects/xyz-finance-go/internal/repository/mock"
	services "github.com/hadi-projects/xyz-finance-go/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLimitService_CreateLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLimitRepo := mock.NewMockLimitRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	service := services.NewLimitService(mockLimitRepo, mockUserRepo)

	t.Run("Success", func(t *testing.T) {
		req := dto.CreateLimitRequest{
			TargetUserID: 1,
			TenorMonth:   12,
			LimitAmount:  1000000,
		}

		limitMatcher := gomock.AssignableToTypeOf(&entity.TenorLimit{})

		mockLimitRepo.EXPECT().
			Create(limitMatcher).
			Do(func(l *entity.TenorLimit) {
				l.ID = 123
			}).
			Return(nil)

		mockUserRepo.EXPECT().
			CreateUserHasTenorLimit(req.TargetUserID, uint(123)).
			Return(nil)

		err := service.CreateLimit(req)
		assert.NoError(t, err)
	})

	t.Run("LimitRepoError", func(t *testing.T) {
		req := dto.CreateLimitRequest{TargetUserID: 1, TenorMonth: 12, LimitAmount: 100}
		mockLimitRepo.EXPECT().Create(gomock.Any()).Return(errors.New("db error"))

		err := service.CreateLimit(req)
		assert.Error(t, err)
		assert.Equal(t, "db error", err.Error())
	})
}

func TestLimitService_DeleteLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLimitRepo := mock.NewMockLimitRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	service := services.NewLimitService(mockLimitRepo, mockUserRepo)

	t.Run("Success", func(t *testing.T) {
		limitID := uint(10)
		mockLimitRepo.EXPECT().Delete(limitID).Return(nil)

		err := service.DeleteLimit(limitID)
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		limitID := uint(10)
		mockLimitRepo.EXPECT().Delete(limitID).Return(errors.New("delete failed"))

		err := service.DeleteLimit(limitID)
		assert.Error(t, err)
	})
}
