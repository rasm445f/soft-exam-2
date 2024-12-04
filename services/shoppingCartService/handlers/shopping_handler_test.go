package handlers

// import (
// 	"bytes"
// 	"encoding/json"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
//
// 	"github.com/go-redis/redismock/v8"
// 	"github.com/rasm445f/soft-exam-2/db"
// 	"github.com/rasm445f/soft-exam-2/domain"
// )
//
// func TestAddItem(t *testing.T) {
// 	// client, mock := redismock.NewClientMock()
// 	client, _ := redismock.NewClientMock()
// 	commands := db.NewShoppingCartRepository(client)
// 	shopping_domain := domain.NewShoppingCartDomain(commands)
// 	// handler := NewShoppingCartHandler(shopping_domain)
// 	mux := http.NewServeMux()
// 	handler := NewShoppingCartHandler(shopping_domain)
// 	mux.Handle("/api/shopping", handler.AddItem())
//
// 	itemParams := domain.AddItemParams{
// 		CustomerId:   123,
// 		RestaurantId: 456,
// 		Name:         "Sample Item",
// 		Price:        30.0,
// 		Quantity:     2,
// 	}
//
// 	itemParamsJSON, err := json.Marshal(itemParams)
// 	if err != nil {
// 		t.Fatalf("unexpected error marshalling item: %v", err)
// 	}
//
// 	req := httptest.NewRequest(http.MethodPost, "/api/shopping", bytes.NewBuffer(itemParamsJSON))
// 	rec := httptest.NewRecorder()
//
// 	mux.ServeHTTP(rec, req)
// 	// handler.AddItem()(rec, req)
//
// 	if rec.Code != http.StatusCreated {
// 		t.Fatalf("got status %d, want %d", rec.Code, http.StatusCreated)
// 	}
//
// }
