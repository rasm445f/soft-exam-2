package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	// "github.com/jackc/pgx/v5/pgconn"
	"github.com/rasm445f/soft-exam-2/db/generated"
	"github.com/rasm445f/soft-exam-2/domain"

	// "github.com/stretchr/testify/assert"
	"github.com/pashagolub/pgxmock/v4"
)

func SetupTestMocks(t *testing.T) (pgxmock.PgxPoolIface, *RestaurantHandler) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create pgxmock pool: %v", err)
	}
	
	queries := generated.New(mock)
	restaurantDomain := domain.NewRestaurantDomain(queries)
	handler := NewRestaurantHandler(restaurantDomain)

	return mock, handler
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

func TestGetAllRestaurantsHandler(t *testing.T) {
	mock, handler := SetupTestMocks(t)
	defer CloseMocks(mock)

	t.Run("Valid Data", func(t *testing.T) {
		// Arrange mock data
		rows := pgxmock.NewRows([]string{"id", "name", "rating", "category", "address", "zip_code"}).
			AddRow(int32(1), "Pizza Paradise", float64Ptr(4.5), stringPtr("Pizza"), stringPtr("Main Street 123"), int32Ptr(2800)).
			AddRow(int32(2), "Sushi World", float64Ptr(4.8), stringPtr("Sushi"), stringPtr("Second Street 456"), int32Ptr(2900))

		mock.ExpectQuery(`SELECT\s+r\.id,\s+r\.name,\s+r\.rating,\s+r\.category,\s+r\.address,\s+r\.zip_code\s+FROM\s+restaurant\s+r\s+JOIN\s+zipcode\s+a\s+ON\s+r\.zip_code\s+=\s+a\.zip_code`).
			WillReturnRows(rows)

		req := httptest.NewRequest(http.MethodGet, "/api/restaurants", nil)
		rec := httptest.NewRecorder()

		// Act
		handler.GetAllRestaurants()(rec, req)

		// Assert
		want := []generated.Restaurant{
			{ID: 1, Name: "Pizza Paradise", Rating: float64Ptr(4.5), Category: stringPtr("Pizza"), Address: stringPtr("Main Street 123"), ZipCode: int32Ptr(2800)},
			{ID: 2, Name: "Sushi World", Rating: float64Ptr(4.8), Category: stringPtr("Sushi"), Address: stringPtr("Second Street 456"), ZipCode: int32Ptr(2900)},
		}

		var got []generated.Restaurant
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %+v, want %+v", got, want)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet mock expectations: %v", err)
		}
	})

	t.Run("Database Error", func(t *testing.T) {
		// Arrange
		mock.ExpectQuery(`SELECT\s+r\.id,\s+r\.name,\s+r\.rating,\s+r\.category,\s+r\.address,\s+r\.zip_code\s+FROM\s+restaurant\s+r\s+JOIN\s+zipcode\s+a\s+ON\s+r\.zip_code\s+=\s+a\.zip_code`).
			WillReturnError(context.DeadlineExceeded)

		req := httptest.NewRequest(http.MethodGet, "/api/restaurants", nil)
		rec := httptest.NewRecorder()
	
		// Act
		handler.GetAllRestaurants()(rec, req)

		// Assert
		if rec.Code != http.StatusInternalServerError {
			t.Fatalf("got status %d, want %d", rec.Code, http.StatusInternalServerError)
		}
		expectedBody := "Failed to get restaurants\n"
		if rec.Body.String() != expectedBody {
			t.Errorf("got body %q, want %q", rec.Body.String(), expectedBody)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet mock expectations: %v", err)
		}
	})

	t.Run("Correct Endpoint Validation", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/wrong-endpoint", nil)
		rec := httptest.NewRecorder()

		// Act
		handler.GetAllRestaurants()(rec, req)

		// Assert
		if rec.Code == http.StatusOK {
			t.Fatalf("expected failure for incorrect endpoint, but got status %d", rec.Code)
		}
	})
}
