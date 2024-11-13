// package handlers_test

// import (
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/jackc/pgx/v5/pgconn"
// 	"github.com/jackc/pgx/v5/pgxpool"
// 	"github.com/rasm445f/soft-exam-2/db/generated"
// 	"github.com/rasm445f/soft-exam-2/handlers"
// )


// // MockPGXPool wraps sqlmock to implement pgx interfaces
// type MockPGXPool struct {
// 	sqlmock.Sqlmock
// }

// func (m MockPGXPool) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
// 	result, err := m.Sqlmock.Exec(query, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return pgconn.CommandTag(result.RowsAffected()), nil
// }

// func (m MockPGXPool) Exec(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
// 	rows, err := m.Sqlmock.Exec(query, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return rows, nil
// }



// // Helper function to create a mock datrabase connection and sqlmock
// func setupMockDB(t *testing.T) (*MockPGXPool, *generated.Queries) {
// 	db, mock, err := sqlmock.New()
// 	if err != nil {
// 		t.Fatalf("failed to create sqlmock: %v", err)
// 	}

// 	pgxMock := &MockPGXPool{Sqlmock: mock}
// 	queries := generated.New(pgxMock)
// 	return pgxMock, queries
// }

// // Helper function to execute a handler and return the response
// func executeRequest(t *testing.T, handler http.HandlerFunc, method, path string) *httptest.ResponseRecorder {
// 	req, err := http.NewRequest(method, path, nil)
// 	if err != nil {
// 		t.Fatalf("failed to create request: %v", err)
// 	}
// 	rec := httptest.NewRecorder()
// 	handler.ServeHTTP(rec, req)
// 	return rec
// }

// // Helper function to assert HTTP status and JSON decoding
// func assertResponse(t *testing.T, rec *httptest.ResponseRecorder, expectedStatus int, target interface{}) {
// 	if rec.Code != expectedStatus {
// 		t.Errorf("exepected status %d, got %d", expectedStatus, rec.Code)
// 	}
// 	if target != nil {
// 		err := json.Unmarshal(rec.Body.Bytes(), target)
// 		if err != nil {
// 			t.Fatalf("failed to unmarshal response: %v", err)
// 		}
// 	}
// }

// // UNIT TESTS
// func TestGetAllRestaurants(t *testing.T) {
// 	mockDB, queries := setupMockDB(t)
// 	defer mockDB.Close()

// 	// Mock expected database query
// 	rows := sqlmock.NewRows([]string{"id", "name", "address", "rating"}).
// 		AddRow(1, "Pizza Paradise", "123 Main Street", 4.5).
// 		AddRow(2, "Burger Bonanza", "456 High Street", 4.2)
// 	mock.ExpectQuery("SELECT id, name, address, rating FROM restaurant").WillReturnRows(rows)

// 	// Execute handler
// 	handler := handlers.GetAllRestaurants(queries)
// 	rec := executeRequest(t, handler, "GET", "/api/restaurants")

// 	// Assert Response
// 	var restaurants []generated.Restaurant
// 	assertResponse(t, rec, http.StatusOK, &restaurants)

// 	if len(restaurants) != 2 {
// 		t.Errorf("expected 2 restaurants, got %d", len(restaurants))
// 	}
// 	if err := mock.ExpectationsWereMet(); err != nil {
// 		t.Errorf("unmet expectations: %v", err)
// 	}
// }