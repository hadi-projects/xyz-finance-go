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
			TenorMonth:   12, // Note: This will fail now due to validation
			LimitAmount:  1000000,
		}
		// Adjust test data for success
		req.TenorMonth = 1

		mockUserRepo.EXPECT().FindByID(req.TargetUserID).Return(&entity.User{ID: 1}, nil)

		// Expect check for existing limits (return empty for success)
		mockLimitRepo.EXPECT().FindByUserID(req.TargetUserID).Return([]entity.TenorLimit{}, nil)

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

	t.Run("InvalidTenor", func(t *testing.T) {
		req := dto.CreateLimitRequest{TargetUserID: 1, TenorMonth: 5, LimitAmount: 100}
		err := service.CreateLimit(req)
		assert.Error(t, err)
		assert.Equal(t, "invalid tenor month: must be 1, 2, 3, or 6", err.Error())
	})

	t.Run("UserNotFound", func(t *testing.T) {
		req := dto.CreateLimitRequest{TargetUserID: 99, TenorMonth: 1, LimitAmount: 100}
		mockUserRepo.EXPECT().FindByID(req.TargetUserID).Return(nil, errors.New("record not found"))

		err := service.CreateLimit(req)
		assert.Error(t, err)
		assert.Equal(t, "target user not found", err.Error())
	})

	t.Run("DuplicateLimit", func(t *testing.T) {
		req := dto.CreateLimitRequest{TargetUserID: 1, TenorMonth: 1, LimitAmount: 100}

		mockUserRepo.EXPECT().FindByID(req.TargetUserID).Return(&entity.User{ID: 1}, nil)

		// Return existing limit for Tenor 1
		mockLimitRepo.EXPECT().FindByUserID(req.TargetUserID).Return([]entity.TenorLimit{
			{TenorMonth: 1, LimitAmount: 50000},
		}, nil)

		err := service.CreateLimit(req)
		assert.Error(t, err)
		assert.Equal(t, "limit for this tenor already exists", err.Error())
	})

	t.Run("LimitRepoError", func(t *testing.T) {
		req := dto.CreateLimitRequest{TargetUserID: 1, TenorMonth: 1, LimitAmount: 100}

		mockUserRepo.EXPECT().FindByID(req.TargetUserID).Return(&entity.User{ID: 1}, nil)
		mockLimitRepo.EXPECT().FindByUserID(req.TargetUserID).Return([]entity.TenorLimit{}, nil)
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
		mockLimitRepo.EXPECT().FindByID(limitID).Return(&entity.TenorLimit{ID: 10}, nil)
		mockLimitRepo.EXPECT().Delete(limitID).Return(nil)

		err := service.DeleteLimit(limitID)
		assert.NoError(t, err)
	})

	t.Run("NotFound", func(t *testing.T) {
		limitID := uint(10)
		mockLimitRepo.EXPECT().FindByID(limitID).Return(nil, errors.New("record not found"))

		err := service.DeleteLimit(limitID)
		assert.Error(t, err)
		assert.Equal(t, "limit not found", err.Error())
	})

	t.Run("DeleteError", func(t *testing.T) {
		limitID := uint(10)
		mockLimitRepo.EXPECT().FindByID(limitID).Return(&entity.TenorLimit{ID: 10}, nil)
		mockLimitRepo.EXPECT().Delete(limitID).Return(errors.New("delete failed"))

		err := service.DeleteLimit(limitID)
		assert.Error(t, err)
		assert.Equal(t, "delete failed", err.Error())
	})
}

func TestLimitService_UpdateLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLimitRepo := mock.NewMockLimitRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	service := services.NewLimitService(mockLimitRepo, mockUserRepo)

	t.Run("Success", func(t *testing.T) {
		limitID := uint(1)
		req := dto.UpdateLimitRequest{
			TenorMonth:  1,
			LimitAmount: 200000,
		}

		mockLimitRepo.EXPECT().FindByID(limitID).Return(&entity.TenorLimit{ID: 1, TenorMonth: 2, LimitAmount: 100000}, nil)
		mockLimitRepo.EXPECT().Update(gomock.Any()).DoAndReturn(func(l *entity.TenorLimit) error {
			assert.Equal(t, entity.Tenor(1), l.TenorMonth)
			assert.Equal(t, 200000.0, l.LimitAmount)
			return nil
		})

		err := service.UpdateLimit(limitID, req)
		assert.NoError(t, err)
	})

	t.Run("InvalidTenor", func(t *testing.T) {
		limitID := uint(1)
		req := dto.UpdateLimitRequest{TenorMonth: 5, LimitAmount: 200000}
		err := service.UpdateLimit(limitID, req)
		assert.Error(t, err)
		assert.Equal(t, "invalid tenor month: must be 1, 2, 3, or 6", err.Error())
	})

	t.Run("NotFound", func(t *testing.T) {
		limitID := uint(1)
		req := dto.UpdateLimitRequest{TenorMonth: 1, LimitAmount: 200000}
		mockLimitRepo.EXPECT().FindByID(limitID).Return(nil, errors.New("record not found"))

		err := service.UpdateLimit(limitID, req)
		assert.Error(t, err)
		assert.Equal(t, "limit not found", err.Error())
	})
}
