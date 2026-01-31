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

func TestTransactionService_CreateTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLimitRepo := mock.NewMockLimitRepository(ctrl)
	mockTxRepo := mock.NewMockTransactionRepository(ctrl)

	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm conn: %v", err)
	}

	service := services.NewTransactionService(mockTxRepo, mockLimitRepo, gormDB)

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
		userId := uint(1)

		sqlMock.ExpectBegin()
		// Expect row locking
		sqlMock.ExpectExec("SELECT id FROM users WHERE id = \\? FOR UPDATE").
			WithArgs(userId).
			WillReturnResult(sqlmock.NewResult(1, 1)) // Return success

		mockLimitRepo.EXPECT().WithTx(gomock.Any()).Return(mockLimitRepo)
		mockTxRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxRepo)

		mockLimitRepo.EXPECT().FindByUserID(userId).Return([]entity.TenorLimit{
			{TenorMonth: 1, LimitAmount: 20000},
		}, nil)

		mockTxRepo.EXPECT().FindByUserID(userId).Return([]entity.Transaction{}, nil)

		sqlMock.ExpectExec("INSERT INTO `transactions`").
			WillReturnResult(sqlmock.NewResult(1, 1))

		sqlMock.ExpectCommit()

		err := service.CreateTransaction(userId, req)
		assert.NoError(t, err)

		if err := sqlMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("InsufficientLimit", func(t *testing.T) {
		req := dto.CreateTransactionRequest{
			OTR:   30000,
			Tenor: 1,
		}
		userId := uint(1)

		sqlMock.ExpectBegin()
		// Expect row locking
		sqlMock.ExpectExec("SELECT id FROM users WHERE id = \\? FOR UPDATE").
			WithArgs(userId).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mockLimitRepo.EXPECT().WithTx(gomock.Any()).Return(mockLimitRepo)
		mockTxRepo.EXPECT().WithTx(gomock.Any()).Return(mockTxRepo)

		mockLimitRepo.EXPECT().FindByUserID(userId).Return([]entity.TenorLimit{
			{TenorMonth: 1, LimitAmount: 20000},
		}, nil)

		mockTxRepo.EXPECT().FindByUserID(userId).Return([]entity.Transaction{}, nil)

		sqlMock.ExpectRollback()

		err := service.CreateTransaction(userId, req)
		assert.Error(t, err)
		assert.Equal(t, "insufficient limit", err.Error())
	})
}
