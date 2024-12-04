package domain

import (
	"context"
	"reflect"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/rasm445f/soft-exam-2/db/generated"
)

func SetupTestMocks(t *testing.T) (pgxmock.PgxPoolIface, *generated.Queries, *RestaurantDomain) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create pgxmock pool: %v", err)
	}
	
	queries := generated.New(mock)
	domain := NewRestaurantDomain(queries)

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

func TestGetAllRestaurantsDomain(t *testing.T) {
	mock, _, domain := SetupTestMocks(t)
	defer CloseMocks(mock)

	t.Run("Valid Data", func(t *testing.T) {
		// Arrange mock data
		rows := pgxmock.NewRows([]string{"id", "name", "rating", "category", "address", "zip_code"}).
			AddRow(int32(1), "Pizza Paradise", float64Ptr(4.5), stringPtr("Pizza"), stringPtr("Main Street 123"), int32Ptr(2800)).
			AddRow(int32(2), "Sushi World", float64Ptr(4.8), stringPtr("Sushi"), stringPtr("Second Street 456"), int32Ptr(2900))

		mock.ExpectQuery(`SELECT\s+r\.id,\s+r\.name,\s+r\.rating,\s+r\.category,\s+r\.address,\s+r\.zip_code\s+FROM\s+restaurant\s+r\s+JOIN\s+zipcode\s+a\s+ON\s+r\.zip_code\s+=\s+a\.zip_code`).
			WillReturnRows(rows)
		
		// Act
		got, err := domain.GetAllRestaurantsDomain(context.Background())

		// Assert
		want := []generated.Restaurant{
			{ID: 1, Name: "Pizza Paradise", Rating: float64Ptr(4.5), Category: stringPtr("Pizza"), Address: stringPtr("Main Street 123"), ZipCode: int32Ptr(2800)},
			{ID: 2, Name: "Sushi World", Rating: float64Ptr(4.8), Category: stringPtr("Sushi"), Address: stringPtr("Second Street 456"), ZipCode: int32Ptr(2900)},
		}

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet mock expectations: %v", err)
		}
	})

	t.Run("No Data", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"id", "name", "rating", "category", "address", "zip_code"})
		mock.ExpectQuery(`SELECT\s+r\.id,\s+r\.name,\s+r\.rating,\s+r\.category,\s+r\.address,\s+r\.zip_code\s+FROM\s+restaurant\s+r\s+JOIN\s+zipcode\s+a\s+ON\s+r\.zip_code\s+=\s+a\.zip_code`).
		WillReturnRows(rows)

		// Act
		got, err := domain.GetAllRestaurantsDomain(context.Background())

		// Assert
		want := []generated.Restaurant{}

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if (got == nil && len(want) != 0) || (got != nil && !reflect.DeepEqual(got, want)) {
			t.Errorf("got %+v, want %+v", got, want)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet mock expectations: %v", err)
		}
	})

	t.Run("Database Error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT\s+r\.id,\s+r\.name,\s+r\.rating,\s+r\.category,\s+r\.address,\s+r\.zip_code\s+FROM\s+restaurant\s+r\s+JOIN\s+zipcode\s+a\s+ON\s+r\.zip_code\s+=\s+a\.zip_code`).
		WillReturnError(context.DeadlineExceeded)

		// Act
		got, err := domain.GetAllRestaurantsDomain(context.Background())

		if err == nil {
			t.Fatalf("expected an error but got nil")
		}
		if got != nil {
			t.Errorf("expected nil but got %+v", got)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet mock expectations: %v", err)
		}
	})
}