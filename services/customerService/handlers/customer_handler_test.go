package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rasm445f/soft-exam-2/db/generated"
)

// Mock the domain layer to test the handler in isolation
type MockCustomerDomain struct {
	GetAllCustomersDomainFunc func(ctx context.Context) ([]generated.GetAllCustomersRow, error)
	GetCustomerByIdDomainFunc func(ctx context.Context, id int32) (generated.GetCustomerByIDRow, error)
	DeleteCustomerDomainFunc  func(ctx context.Context, id int32) error
	CreateCustomerDomainFunc  func(ctx context.Context, customerParams generated.CreateCustomerParams) error
	UpdateCustomerDomainFunc  func(ctx context.Context, customerParams generated.UpdateCustomerParams) error
	UpdateAddressFunc         func(ctx context.Context, addressParams generated.UpdateAddressParams) error
}

func int32Ptr(i int32) *int32 {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func (m *MockCustomerDomain) GetAllCustomersDomain(ctx context.Context) ([]generated.GetAllCustomersRow, error) {
	if m.GetAllCustomersDomainFunc != nil {
		return m.GetAllCustomersDomainFunc(ctx)
	}
	return nil, nil
}

func (m *MockCustomerDomain) GetCustomerByIdDomain(ctx context.Context, id int32) (generated.GetCustomerByIDRow, error) {
	if m.GetCustomerByIdDomainFunc != nil {
		return m.GetCustomerByIdDomainFunc(ctx, id)
	}
	return generated.GetCustomerByIDRow{}, nil
}

func (m *MockCustomerDomain) DeleteCustomerDomain(ctx context.Context, id int32) error {
	if m.DeleteCustomerDomainFunc != nil {
		return m.DeleteCustomerDomainFunc(ctx, id)
	}
	return nil
}

func (m *MockCustomerDomain) CreateCustomerDomain(ctx context.Context, customerParams generated.CreateCustomerParams) error {
	if m.CreateCustomerDomainFunc != nil {
		return m.CreateCustomerDomainFunc(ctx, customerParams)
	}
	return nil
}

func (m *MockCustomerDomain) UpdateCustomerDomain(ctx context.Context, customerParams generated.UpdateCustomerParams) error {
	if m.UpdateCustomerDomainFunc != nil {
		return m.UpdateCustomerDomainFunc(ctx, customerParams)
	}
	return nil
}

func (m *MockCustomerDomain) UpdateAddress(ctx context.Context, addressParams generated.UpdateAddressParams) error {
	if m.UpdateAddressFunc != nil {
		return m.UpdateAddressFunc(ctx, addressParams)
	}
	return nil
}

func TestGetAllCustomersHandler(t *testing.T) {
	mockDomain := &MockCustomerDomain{
		GetAllCustomersDomainFunc: func(ctx context.Context) ([]generated.GetAllCustomersRow, error) {
			return []generated.GetAllCustomersRow{
				{ID: 1, Name: stringPtr("John Doe"), Email: stringPtr("john@example.com")},
				{ID: 2, Name: stringPtr("Jane Smith"), Email: stringPtr("jane@example.com")},
			}, nil
		},
	}
	handler := NewCustomerHandler(mockDomain)

	t.Run("status 200", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/customers", nil)
		handler.GetAllCustomers().ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusOK {
			t.Fatalf("expected status %v, got %v", http.StatusOK, rec.Result().StatusCode)
		}
	})

	t.Run("internal server error", func(t *testing.T) {
		mockDomain.GetAllCustomersDomainFunc = func(ctx context.Context) ([]generated.GetAllCustomersRow, error) {
			return nil, sql.ErrConnDone
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/customers", nil)
		handler.GetAllCustomers().ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusInternalServerError {
			t.Fatalf("expected status %v, got %v", http.StatusInternalServerError, rec.Result().StatusCode)
		}
	})
}

func TestCreateCustomerHandler(t *testing.T) {
	mockDomain := &MockCustomerDomain{
		CreateCustomerDomainFunc: func(ctx context.Context, customerParams generated.CreateCustomerParams) error {
			return nil
		},
	}
	handler := NewCustomerHandler(mockDomain)

	t.Run("status 201", func(t *testing.T) {
		customer := generated.CreateCustomerParams{
			Name:          stringPtr("John Doe"),
			Email:         stringPtr("john@example.com"),
			Phonenumber:   stringPtr("1234567890"),
			Password:      stringPtr("Password123!"),
			StreetAddress: stringPtr("123 Main St"),
			ZipCode:       int32Ptr(12345),
		}

		customerJSON, _ := json.Marshal(customer)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/customer", bytes.NewBuffer(customerJSON))
		handler.CreateCustomer().ServeHTTP(rec, req)

		got := rec.Result().StatusCode
		want := http.StatusCreated
		// Assertions
		if rec.Result().StatusCode != http.StatusCreated {
			t.Fatalf("expected status %v, got %v", want, got)
		}
	})

	t.Run("invalid email", func(t *testing.T) {
		params := generated.CreateCustomerParams{
			StreetAddress: stringPtr("123 Main St"),
			ZipCode:       int32Ptr(12345),
			Name:          stringPtr("John Doe"),
			Email:         stringPtr(""),
			Phonenumber:   stringPtr("12341212"),
			Password:      stringPtr("Password1!"),
		}
		paramsJSON, _ := json.Marshal(params)
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewBuffer(paramsJSON))
		handler.CreateCustomer().ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status %v, got %v", http.StatusBadRequest, rec.Result().StatusCode)
		}
	})

	t.Run("invalid payload", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/customer", bytes.NewBuffer([]byte(`{invalid json}`)))
		handler.CreateCustomer().ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status %v, got %v", http.StatusBadRequest, rec.Result().StatusCode)
		}
	})

	t.Run("internal server error", func(t *testing.T) {
		mockDomain.CreateCustomerDomainFunc = func(ctx context.Context, customerParams generated.CreateCustomerParams) error {
			return sql.ErrConnDone
		}

		customer := generated.CreateCustomerParams{
			Name:          stringPtr("John Doe"),
			Email:         stringPtr("john@example.com"),
			Phonenumber:   stringPtr("1234567890"),
			Password:      stringPtr("Password123!"),
			StreetAddress: stringPtr("123 Main St"),
			ZipCode:       int32Ptr(12345),
		}
		customerJSON, _ := json.Marshal(customer)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/customer", bytes.NewBuffer(customerJSON))
		handler.CreateCustomer().ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusInternalServerError {
			t.Fatalf("expected status %v, got %v", http.StatusInternalServerError, rec.Result().StatusCode)
		}
	})
}

func TestGetCustomerByIdHandler(t *testing.T) {
	mockDomain := &MockCustomerDomain{
		GetCustomerByIdDomainFunc: func(ctx context.Context, id int32) (generated.GetCustomerByIDRow, error) {
			if id == 1 {
				return generated.GetCustomerByIDRow{
					ID:            1,
					Name:          stringPtr("John Doe"),
					Email:         stringPtr("charlie@example.com"),
					Phonenumber:   stringPtr("12341212"),
					Password:      stringPtr("Password1!"),
					StreetAddress: stringPtr("123 Main St"),
					ZipCode:       int32Ptr(12345),
					City:          stringPtr("New York"),
				}, nil
			}
			return generated.GetCustomerByIDRow{}, sql.ErrNoRows
		},
	}
	handler := NewCustomerHandler(mockDomain)

	t.Run("status 200", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/customer/1", nil)
		req.SetPathValue("id", "1")

		handler.GetCustomerById().ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusOK {
			t.Fatalf("expected status %v, got %v", http.StatusOK, rec.Result().StatusCode)
		}
	})

	t.Run("not found", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/customer/999", nil)
		req.SetPathValue("id", "999")

		handler.GetCustomerById().ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusNotFound {
			t.Fatalf("expected status %v, got %v", http.StatusNotFound, rec.Result().StatusCode)
		}
	})

	t.Run("invalid ID", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/customer/invalid", nil)
		req.SetPathValue("id", "invalid")

		handler.GetCustomerById().ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status %v, got %v", http.StatusBadRequest, rec.Result().StatusCode)
		}
	})
}

func TestUpdateCustomerHandler(t *testing.T) {
	mockDomain := &MockCustomerDomain{
		UpdateCustomerDomainFunc: func(ctx context.Context, customerParams generated.UpdateCustomerParams) error {
			return nil
		},
		UpdateAddressFunc: func(ctx context.Context, addressParams generated.UpdateAddressParams) error {
			return nil
		},
	}
	handler := NewCustomerHandler(mockDomain)

	t.Run("status 200", func(t *testing.T) {
		customerUpdates := UpdateCustomerWithAddress{
			Name:          stringPtr("John Doe"),
			Email:         stringPtr("john@example.com"),
			Phonenumber:   stringPtr("1234567890"),
			Password:      stringPtr("Password123!"),
			StreetAddress: stringPtr("123 Main St"),
			ZipCode:       int32Ptr(12345),
			City:          stringPtr("New York"),
		}
		customerUpdatesJSON, _ := json.Marshal(customerUpdates)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/api/customer/1", bytes.NewBuffer(customerUpdatesJSON))
		req.SetPathValue("id", "1")
		handler.UpdateCustomer().ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusOK {
			t.Fatalf("expected status %v, got %v", http.StatusOK, rec.Result().StatusCode)
		}
	})

	t.Run("invalid payload", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/api/customer/1", bytes.NewBuffer([]byte(`{invalid json}`)))
		req.SetPathValue("id", "1")
		handler.UpdateCustomer().ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status %v, got %v", http.StatusBadRequest, rec.Result().StatusCode)
		}
	})

	t.Run("customer not found", func(t *testing.T) {
		mockDomain.UpdateCustomerDomainFunc = func(ctx context.Context, customerParams generated.UpdateCustomerParams) error {
			return sql.ErrNoRows
		}

		customerUpdates := UpdateCustomerWithAddress{
			Name:          stringPtr("John Doe"),
			Email:         stringPtr("john@example.com"),
			Phonenumber:   stringPtr("1234567890"),
			Password:      stringPtr("Password123!"),
			StreetAddress: stringPtr("123 Main St"),
			ZipCode:       int32Ptr(12345),
			City:          stringPtr("New York"),
		}
		customerUpdatesJSON, _ := json.Marshal(customerUpdates)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/api/customer/999", bytes.NewBuffer(customerUpdatesJSON))
		req.SetPathValue("id", "999")
		handler.UpdateCustomer().ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusNotFound {
			t.Fatalf("expected status %v, got %v", http.StatusNotFound, rec.Result().StatusCode)
		}
	})

	t.Run("invalid ID", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/api/customer/invalid", nil)
		req.SetPathValue("id", "invalid")
		handler.UpdateCustomer().ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status %v, got %v", http.StatusBadRequest, rec.Result().StatusCode)
		}
	})

	// t.Run("failed to update address", func(t *testing.T) {
	// 	mockDomain.UpdateAddressFunc = func(ctx context.Context, addressParams generated.UpdateAddressParams) error {
	// 		return sql.ErrConnDone
	// 	}

	// 	customerUpdates := UpdateCustomerWithAddress{
	// 		Name:          stringPtr("John Doe"),
	// 		Email:         stringPtr("john@example.com"),
	// 		Phonenumber:   stringPtr("1234567890"),
	// 		Password:      stringPtr("Password123!"),
	// 		StreetAddress: stringPtr("123 Main St"),
	// 		ZipCode:       int32Ptr(12345),
	// 	}
	// 	customerUpdatesJSON, _ := json.Marshal(customerUpdates)

	// 	rec := httptest.NewRecorder()
	// 	req := httptest.NewRequest(http.MethodPatch, "/api/customer/1", bytes.NewBuffer(customerUpdatesJSON))
	// 	req.SetPathValue("id", "1")
	// 	handler.UpdateCustomer().ServeHTTP(rec, req)

	// 	if rec.Result().StatusCode != http.StatusInternalServerError {
	// 		t.Fatalf("expected status %v, got %v", http.StatusInternalServerError, rec.Result().StatusCode)
	// 	}
	// })

}

func TestDeleteCustomerHandler(t *testing.T) {
	mockDomain := &MockCustomerDomain{
		DeleteCustomerDomainFunc: func(ctx context.Context, id int32) error {
			if id == 1 {
				return nil
			}
			return sql.ErrNoRows
		},
	}
	handler := NewCustomerHandler(mockDomain)

	t.Run("status 200", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/api/customer/1", nil)
		req.SetPathValue("id", "1")
		handler.DeleteCustomer().ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusOK {
			t.Fatalf("expected status %v, got %v", http.StatusOK, rec.Result().StatusCode)
		}
	})

	t.Run("not found", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/api/customer/999", nil)
		req.SetPathValue("id", "999")
		handler.DeleteCustomer().ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusNotFound {
			t.Fatalf("expected status %v, got %v", http.StatusNotFound, rec.Result().StatusCode)
		}
	})

	t.Run("invalid ID", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/api/customer/invalid", nil)
		req.SetPathValue("id", "invalid")
		handler.DeleteCustomer().ServeHTTP(rec, req)

		if rec.Result().StatusCode != http.StatusBadRequest {
			t.Fatalf("expected status %v, got %v", http.StatusBadRequest, rec.Result().StatusCode)
		}
	})

}
