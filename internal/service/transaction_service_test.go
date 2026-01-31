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
	mockMutationRepo := mock.NewMockLimitMutationRepository(ctrl)

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

	service := services.NewTransactionService(mockTxRepo, mockLimitRepo, mockMutationRepo, gormDB)

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
		mockMutationRepo.EXPECT().WithTx(gomock.Any()).Return(mockMutationRepo)

		mockLimitRepo.EXPECT().FindByUserID(userId).Return([]entity.TenorLimit{
			{ID: 123, TenorMonth: 1, LimitAmount: 20000},
		}, nil)

		mockTxRepo.EXPECT().FindByUserID(userId).Return([]entity.Transaction{}, nil)

		mockTxRepo.EXPECT().Create(gomock.Any()).Return(nil)

		// Expect Mutation Logging
		mockMutationRepo.EXPECT().Create(gomock.Any()).Do(func(m *entity.LimitMutation) {
			assert.Equal(t, entity.MutationUsage, m.Action)
			assert.Equal(t, uint(123), m.TenorLimitID)
			assert.Equal(t, 20000.0, m.OldAmount)
			assert.Equal(t, 20000.0, m.NewAmount)
			assert.Contains(t, m.Reason, "Transaction Usage")
		}).Return(nil)

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
		// Note: Mutation repo might not be called if error happens early, or WithTx might be called.
		// Depending on implementation order. In current impl, WithTx is called early.
		// mockMutationRepo.EXPECT().WithTx(gomock.Any()).Return(mockMutationRepo) // Not called because WithTx is called inside transaction block, and limit check fails inside. Wait.
		// The WithTx calls are at the beginning of Transaction block. So expectations needed.
		// But in the "Success" test above, I saw WithTx being called.
		// Let's re-read the code.
		// Service code:
		// limitRepoTx := s.limitRepo.WithTx(tx)
		// transactionRepoTx := s.transactionRepo.WithTx(tx)
		// ... Limit Check ...
		// If limit check fails, it returns error, triggering Rollback.
		// Mutation Create is NOT called.
		// So we need WithTx expectation, but NO Create expectation.

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
