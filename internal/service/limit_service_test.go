package services_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hadi-projects/xyz-finance-go/internal/dto"
	"github.com/hadi-projects/xyz-finance-go/internal/entity"
	"github.com/hadi-projects/xyz-finance-go/internal/repository/mock"
	services "github.com/hadi-projects/xyz-finance-go/internal/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestLimitService_CreateLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a gorm database connection", err)
	}

	mockLimitRepo := mock.NewMockLimitRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	mockMutationRepo := mock.NewMockLimitMutationRepository(ctrl)

	service := services.NewLimitService(mockLimitRepo, mockUserRepo, mockMutationRepo, gormDB)

	t.Run("Success", func(t *testing.T) {
		req := dto.CreateLimitRequest{
			TargetUserID: 1,
			TenorMonth:   1,
			LimitAmount:  1000000,
		}

		mockUserRepo.EXPECT().FindByID(req.TargetUserID).Return(&entity.User{ID: 1}, nil)
		mockLimitRepo.EXPECT().FindByUserID(req.TargetUserID).Return([]entity.TenorLimit{}, nil)

		// Expect Transaction
		sqlMock.ExpectBegin()

		// Expect WithTx
		mockLimitRepo.EXPECT().WithTx(gomock.Any()).Return(mockLimitRepo)
		mockMutationRepo.EXPECT().WithTx(gomock.Any()).Return(mockMutationRepo)

		// Expect Create Limit
		mockLimitRepo.EXPECT().Create(gomock.Any()).Do(func(l *entity.TenorLimit) {
			l.ID = 123
		}).Return(nil)

		// Expect Manual Association Insert
		sqlMock.ExpectExec("INSERT INTO user_has_tenor_limit").
			WithArgs(req.TargetUserID, 123).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// Expect Create Mutation
		mockMutationRepo.EXPECT().Create(gomock.Any()).Do(func(m *entity.LimitMutation) {
			assert.Equal(t, uint(1), m.UserID)
			assert.Equal(t, uint(123), m.TenorLimitID)
			assert.Equal(t, entity.MutationCreate, m.Action)
			assert.Equal(t, 0.0, m.OldAmount)
			assert.Equal(t, 1000000.0, m.NewAmount)
		}).Return(nil)

		sqlMock.ExpectCommit()

		err := service.CreateLimit(req)
		assert.NoError(t, err)

		if err := sqlMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("DuplicateLimit", func(t *testing.T) {
		req := dto.CreateLimitRequest{TargetUserID: 1, TenorMonth: 1, LimitAmount: 100}
		mockUserRepo.EXPECT().FindByID(req.TargetUserID).Return(&entity.User{ID: 1}, nil)
		mockLimitRepo.EXPECT().FindByUserID(req.TargetUserID).Return([]entity.TenorLimit{
			{TenorMonth: 1, LimitAmount: 50000},
		}, nil)

		err := service.CreateLimit(req)
		assert.Error(t, err)
		assert.Equal(t, "limit for this tenor already exists", err.Error())
	})
}

func TestLimitService_UpdateLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a gorm database connection", err)
	}

	mockLimitRepo := mock.NewMockLimitRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	mockMutationRepo := mock.NewMockLimitMutationRepository(ctrl)
	service := services.NewLimitService(mockLimitRepo, mockUserRepo, mockMutationRepo, gormDB)

	t.Run("Success", func(t *testing.T) {
		limitID := uint(1)
		req := dto.UpdateLimitRequest{
			TenorMonth:  1,
			LimitAmount: 200000,
		}

		sqlMock.ExpectBegin()
		mockLimitRepo.EXPECT().WithTx(gomock.Any()).Return(mockLimitRepo)
		mockMutationRepo.EXPECT().WithTx(gomock.Any()).Return(mockMutationRepo)

		mockLimitRepo.EXPECT().FindByID(limitID).Return(&entity.TenorLimit{ID: 1, TenorMonth: 2, LimitAmount: 100000}, nil)
		mockLimitRepo.EXPECT().GetUserIDByLimitID(limitID).Return(uint(101), nil)

		mockLimitRepo.EXPECT().Update(gomock.Any()).Return(nil)

		mockMutationRepo.EXPECT().Create(gomock.Any()).Do(func(m *entity.LimitMutation) {
			assert.Equal(t, uint(101), m.UserID)
			assert.Equal(t, uint(1), m.TenorLimitID)
			assert.Equal(t, entity.MutationUpdate, m.Action)
			assert.Equal(t, 100000.0, m.OldAmount)
			assert.Equal(t, 200000.0, m.NewAmount)
		}).Return(nil)

		sqlMock.ExpectCommit()

		err := service.UpdateLimit(limitID, req)
		assert.NoError(t, err)
	})
}

func TestLimitService_DeleteLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a gorm database connection", err)
	}

	mockLimitRepo := mock.NewMockLimitRepository(ctrl)
	mockUserRepo := mock.NewMockUserRepository(ctrl)
	mockMutationRepo := mock.NewMockLimitMutationRepository(ctrl)
	service := services.NewLimitService(mockLimitRepo, mockUserRepo, mockMutationRepo, gormDB)

	t.Run("Success", func(t *testing.T) {
		limitID := uint(10)

		sqlMock.ExpectBegin()
		mockLimitRepo.EXPECT().WithTx(gomock.Any()).Return(mockLimitRepo)
		mockMutationRepo.EXPECT().WithTx(gomock.Any()).Return(mockMutationRepo)

		mockLimitRepo.EXPECT().FindByID(limitID).Return(&entity.TenorLimit{ID: 10, LimitAmount: 50000}, nil)
		mockLimitRepo.EXPECT().GetUserIDByLimitID(limitID).Return(uint(101), nil)
		mockLimitRepo.EXPECT().Delete(limitID).Return(nil)

		mockMutationRepo.EXPECT().Create(gomock.Any()).Do(func(m *entity.LimitMutation) {
			assert.Equal(t, uint(101), m.UserID)
			assert.Equal(t, entity.MutationDelete, m.Action)
			assert.Equal(t, 50000.0, m.OldAmount)
			assert.Equal(t, 0.0, m.NewAmount)
		}).Return(nil)

		sqlMock.ExpectCommit()

		err := service.DeleteLimit(limitID)
		assert.NoError(t, err)
	})
}
