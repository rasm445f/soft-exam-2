package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rasm445f/soft-exam-2/db"
	"github.com/rasm445f/soft-exam-2/domain"
)

// Mock the domain layer to test the handler in isolation
type MockShoppingCartDomain struct {
	AddItemDomainFunc    func(ctx context.Context, params domain.AddItemParams) error
	UpdateCartDomainFunc func(ctx context.Context, customerId, itemID, quantity int) error
	ViewCartDomainFunc   func(ctx context.Context, customerId int) (*db.ShoppingCart, error)
	ClearCartDomainFunc  func(ctx context.Context, customerId int) error
}

func (m *MockShoppingCartDomain) AddItemDomain(ctx context.Context, params domain.AddItemParams) error {
	if m.AddItemDomainFunc != nil {
		return m.AddItemDomainFunc(ctx, params)
	}
	return nil
}

func (m *MockShoppingCartDomain) UpdateCartDomain(ctx context.Context, customerId, itemID, quantity int) error {
	if m.UpdateCartDomainFunc != nil {
		return m.UpdateCartDomainFunc(ctx, customerId, itemID, quantity)
	}
	return nil
}

func (m *MockShoppingCartDomain) ViewCartDomain(ctx context.Context, customerId int) (*db.ShoppingCart, error) {
	if m.ViewCartDomainFunc != nil {
		return m.ViewCartDomainFunc(ctx, customerId)
	}
	return nil, nil
}

func (m *MockShoppingCartDomain) ClearCartDomain(ctx context.Context, customerId int) error {
	if m.ClearCartDomainFunc != nil {
		return m.ClearCartDomainFunc(ctx, customerId)
	}
	return nil
}

func TestAddItem(t *testing.T) {
	mockDomain := &MockShoppingCartDomain{
		// Actually dont need to specify since the function is defined above
		// AddItemDomainFunc: func(ctx context.Context, params domain.AddItemParams) error {
		// 	return nil
		// },
	}

	handler := NewShoppingCartHandler(mockDomain)

	t.Run("status 201", func(t *testing.T) {
		item := domain.AddItemParams{
			CustomerId:   123,
			RestaurantId: 456,
			Name:         "ting",
			Price:        1,
			Quantity:     1,
		}

		itemJSON, _ := json.Marshal(item)

		rec := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "", bytes.NewBuffer(itemJSON))

		if err != nil {
			t.Error(err)
		}

		handler.AddItem().ServeHTTP(rec, req)

		got := rec.Result().StatusCode
		want := http.StatusCreated
		// Assertions
		if got != want {
			t.Fatalf("expected status %v, got %v", want, got)
		}
	})

}

func TestUpdateCartHandler(t *testing.T) {
	mockDomain := &MockShoppingCartDomain{
		// UpdateCartDomainFunc: func(ctx context.Context, customerId, itemID, quantity int) error {
		// 	return nil
		// },
	}
	handler := NewShoppingCartHandler(mockDomain)

	// Test data
	updateRequest := UpdateQuantityRequest{
		Quantity: 3,
	}
	updateRequestJSON, err := json.Marshal(updateRequest)
	if err != nil {
		t.Error(err)
	}

	// Create a ResponseRecorder
	rec := httptest.NewRecorder()

	// Create a new HTTP request
	req, err := http.NewRequest(http.MethodPatch, "", bytes.NewBuffer(updateRequestJSON))
	if err != nil {
		t.Error(err)
	}
	req.SetPathValue("customerId", "123")
	req.SetPathValue("itemId", "456")
	req.Header.Set("Content-Type", "application/json")

	// call the handler function with the recorder and request
	handler.UpdateCart().ServeHTTP(rec, req)

	got := rec.Result().StatusCode
	want := http.StatusOK
	// Assertions
	if got != want {
		t.Fatalf("expected status %v, got %v", want, got)
	}
}
