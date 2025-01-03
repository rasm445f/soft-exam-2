package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rasm445f/soft-exam-2/broker"
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
	cart := db.ShoppingCart{
		CustomerId:   customerId,
		RestaurantId: 1,
		TotalAmount:  20,
		VatAmount:    4,
		Items:        []db.ShoppingCartItem{},
	}
	return &cart, nil
}

func (m *MockShoppingCartDomain) ClearCartDomain(ctx context.Context, customerId int) error {
	if m.ClearCartDomainFunc != nil {
		return m.ClearCartDomainFunc(ctx, customerId)
	}
	return nil
}

func TestAddItem(t *testing.T) {
	mockDomain := &MockShoppingCartDomain{}
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

	t.Run("invalid request body", func(t *testing.T) {
		reqBody := `{"itemID":123` // malformed JSON
		r := httptest.NewRequest(http.MethodPost, "/add-item", bytes.NewBuffer([]byte(reqBody)))
		w := httptest.NewRecorder()

		handler.AddItem().ServeHTTP(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %v, got %v", http.StatusBadRequest, w.Code)
		}

		expectedBody := "Invalid request body\n"
		if w.Body.String() != expectedBody {
			t.Errorf("expected body %q, got %q", expectedBody, w.Body.String())
		}
	})

	t.Run("domain error", func(t *testing.T) {
		mockDomain := &MockShoppingCartDomain{
			AddItemDomainFunc: func(ctx context.Context, params domain.AddItemParams) error {
				return errors.New("some domain error")
			},
		}
		handler := NewShoppingCartHandler(mockDomain)

		item := domain.AddItemParams{
			CustomerId:   123,
			RestaurantId: 456,
			Name:         "ting",
			Price:        1,
			Quantity:     1,
		}

		itemJSON, _ := json.Marshal(item)
		r := httptest.NewRequest(http.MethodPost, "/add-item", bytes.NewBuffer([]byte(itemJSON)))
		w := httptest.NewRecorder()

		handler.AddItem().ServeHTTP(w, r)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %v, got %v", http.StatusBadRequest, w.Code)
		}

		expectedBody := "some domain error\n"
		if w.Body.String() != expectedBody {
			t.Errorf("expected body %q, got %q", expectedBody, w.Body.String())
		}
	})
}

func TestUpdateCartHandler(t *testing.T) {
	mockDomain := &MockShoppingCartDomain{}
	handler := NewShoppingCartHandler(mockDomain)

	// Test data
	updateRequest := UpdateQuantityRequest{
		Quantity: 3,
	}
	updateRequestJSON, err := json.Marshal(updateRequest)
	if err != nil {
		t.Error(err)
	}

	t.Run("update cart", func(t *testing.T) {
		// Create a ResponseRecorder
		rec := httptest.NewRecorder()
		// Create a new HTTP request
		// req, _ := http.NewRequest(http.MethodPatch, "", bytes.NewBuffer(updateRequestJSON))
		req := httptest.NewRequest(http.MethodPatch, "/", bytes.NewBuffer(updateRequestJSON))
		req.SetPathValue("customerId", "123")
		req.SetPathValue("itemId", "456")
		req.Header.Set("Content-Type", "application/json")

		// call the handler function with the recorder and request
		handler.UpdateCart().ServeHTTP(rec, req)

		got := rec.Result().StatusCode
		want := http.StatusOK
		// Assertions
		if got != want {
			t.Fatalf("want status %v, got %v", want, got)
		}
	})
	t.Run("malformed path values", func(t *testing.T) {
		// Create a ResponseRecorder
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/", bytes.NewBuffer(updateRequestJSON))
		req.SetPathValue("customerId", "somerandomtext")
		req.SetPathValue("itemId", "moretext")
		req.Header.Set("Content-Type", "application/json")

		// call the handler function with the recorder and request
		handler.UpdateCart().ServeHTTP(rec, req)

		got := rec.Result().StatusCode
		want := http.StatusBadRequest
		// Assertions
		if got != want {
			t.Fatalf("want status %v, got %v", want, got)
		}

	})
	t.Run("malformed json body", func(t *testing.T) {
		reqBody := `{"itemID":123` // malformed JSON
		// Create a ResponseRecorder
		rec := httptest.NewRecorder()
		// Create a new HTTP request
		// req, _ := http.NewRequest(http.MethodPatch, "", bytes.NewBuffer(updateRequestJSON))
		req := httptest.NewRequest(http.MethodPatch, "/", bytes.NewBuffer([]byte(reqBody)))
		req.SetPathValue("customerId", "123")
		req.SetPathValue("itemId", "456")
		req.Header.Set("Content-Type", "application/json")

		// call the handler function with the recorder and request
		handler.UpdateCart().ServeHTTP(rec, req)

		got := rec.Result().StatusCode
		want := http.StatusBadRequest
		// Assertions
		if got != want {
			t.Fatalf("want status %v, got %v", want, got)
		}
	})
}

func TestViewCart(t *testing.T) {
	mockDomain := &MockShoppingCartDomain{}
	handler := NewShoppingCartHandler(mockDomain)
	rec := httptest.NewRecorder()

	// Create a new HTTP request
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	req.SetPathValue("customerId", "123")
	handler.ViewCart().ServeHTTP(rec, req)

	// Assertions
	got := rec.Result().StatusCode
	want := http.StatusOK
	if got != want {
		t.Fatalf("expected status %v, got %v", want, got)
	}

	// TODO: check the body
	body, err := io.ReadAll(rec.Result().Body)
	if err != nil {
		t.Error(err)
	}
	defer rec.Result().Body.Close()

	fmt.Println("ting", string(body))
}

func TestClearCart(t *testing.T) {
	mockDomain := &MockShoppingCartDomain{}
	handler := NewShoppingCartHandler(mockDomain)
	rec := httptest.NewRecorder()

	// Create a new HTTP request
	req, _ := http.NewRequest(http.MethodDelete, "", nil)
	req.SetPathValue("customerId", "123")
	handler.ClearCart().ServeHTTP(rec, req)

	// Assertions
	got := rec.Result().StatusCode
	want := http.StatusOK
	if got != want {
		t.Fatalf("expected status %v, got %v", want, got)
	}
}

func TestConsumeMenuItem(t *testing.T) {
	mockDomain := &MockShoppingCartDomain{}
	handler := NewShoppingCartHandler(mockDomain)
	// setup broker and simulate published message to consume
	broker.InitRabbitMQ()
	defer broker.CloseRabbitMQ()

	menuItemSelection := domain.AddItemParams{
		CustomerId:   1,
		RestaurantId: 1,
		Name:         "pizza",
		Price:        20,
		Quantity:     1,
	}

	event := broker.Event{
		Type:    broker.MenuItemSelected,
		Payload: menuItemSelection,
	}
	err := broker.Publish("menu_item_selected_queue", event)
	if err != nil {
		t.Error(err)
	}

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	handler.ConsumeMenuItem().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("got status %d, want %d", rec.Code, http.StatusOK)
	}
	got, _ := io.ReadAll(rec.Body)
	want := `{"message": "Menu item added to cart successfully"}`
	if string(got) != want {
		t.Errorf("expected body %q, got %q", want, string(got))
	}

}
