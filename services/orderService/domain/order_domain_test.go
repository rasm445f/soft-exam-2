package domain

import (
	"context"
	// "errors"
	// "reflect"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/rasm445f/soft-exam-2/db/generated"
)

func SetupTestMocks(t *testing.T) (pgxmock.PgxPoolIface, *generated.Queries, *OrderDomain) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create pgxmock pool: %v", err)
	}
		
	queries := generated.New(mock)
	domain := NewOrderDomain(queries)
	
	return mock, queries, domain
}
	
func CloseMocks(mock pgxmock.PgxPoolIface) {
	mock.Close()
}

// Helper functions to create pointers for literals
func int32Ptr(i int32) *int32 {
	return &i
}
func float64Ptr(f float64) *float64 {
	return &f
}
func stringPtr(s string) *string {
	return &s
}

func TestOrderCalculationLogicDomain(t *testing.T) {
	mock, _, domain := SetupTestMocks(t)
	defer CloseMocks(mock)

		t.Run("Valid Fee Calculation", func(t *testing.T) {
			// Arrange
			amount := float64Ptr(200.0)
			expectedFee := float64Ptr(10.0)
			expectedPercentage := float64Ptr(0.05)
			mock.ExpectQuery(`INSERT INTO Fee`).
				WithArgs(expectedPercentage, expectedFee, stringPtr("some description")).
				WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(int32(1)))
			
			// Act
			feeId, err := domain.CalculateFee(context.Background(), *amount)

			// Assert
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if feeId != 1 {
				t.Errorf("got fee ID %d, want 1", feeId)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("unmet mock expectations: %v", err)
			}
		})

		// t.Run("Database Error", func(t *testing.T) {

		// 	// Arange
		// 	amount := 200.0
		// 	mock.ExpectQuery(`INSERT INTO Fee`).
		// 		WillReturnError(errors.New("database error"))

		// 	// Act
		// 	feeId, err := domain.CalculateFee(context.Background(), amount)

		// 	// Assert
		// 	if err == nil {
		// 		t.Fatalf("expected")
		// 	}
		// })
}