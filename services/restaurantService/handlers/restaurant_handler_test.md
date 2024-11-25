// package handlers

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"math/big"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/jackc/pgx/v5/pgtype"
// 	"github.com/rasm445f/soft-exam-2/db/generated"
// 	"github.com/stretchr/testify/mock"
// )

// // MockRestaurantDomain mocks the RestaurantDomain for testing handlers
// type MockRestaurantDomain struct {
// 	mock.Mock
// }

// func (m *MockRestaurantDomain) FetchAllRestaurants(ctx context.Context) ([]generated.Restaurant, error) {
// 	args := m.Called(ctx)
// 	return args.Get(0).([]generated.Restaurant), args.Error(1)
// }

// func (m *MockRestaurantDomain) GetRestaurantById(ctx context.Context, restaurantId int32) (*generated.Restaurant, error) {
// 	args := m.Called(ctx, restaurantId)
// 	return args.Get(0).(*generated.Restaurant), args.Error(1)
// }

// func (m *MockRestaurantDomain) FetchMenuItemsByRestaurantId(ctx context.Context, restaurantId int32) ([]generated.FetchMenuItemsByRestaurantIdRow, error) {
// 	args := m.Called(ctx, restaurantId)
// 	return args.Get(0).([]generated.FetchMenuItemsByRestaurantIdRow), args.Error(1)
// }

// func (m *MockRestaurantDomain) GetMenuItemByRestaurantAndId(ctx context.Context, params generated.GetMenuItemByRestaurantAndIdParams) (*generated.Menuitem, error) {
// 	args := m.Called(ctx, params)
// 	return args.Get(0).(*generated.Menuitem), args.Error(1)
// }

// func (m *MockRestaurantDomain) FetchAllCategories(ctx context.Context) ([]string, error) {
// 	args := m.Called(ctx)
// 	return args.Get(0).([]string), args.Error(1)
// }

// func (m *MockRestaurantDomain) FilterRestaurantsByCategory(ctx context.Context, category string) ([]generated.Restaurant, error) {
// 	args := m.Called(ctx, category)
// 	return args.Get(0).([]generated.Restaurant), args.Error(1)
// }

// func ptr(s string) *string {
// 	return &s
// }

// func bigIntFromFloat(f float64) *big.Int {
// 	bigInt := new(big.Int)
// 	bigInt.SetString(fmt.Sprintf("%.0f", f*10), 10)
// 	return bigInt
// }

// /* UNIT TESTS */
// func TestGetAllRestaurants(t *testing.T) {
// 	mockDomain := new(MockRestaurantDomain)
// 	handler := NewRestaurantHandler(mockDomain)

// 	// Mock response
// 	expectedRestaurants := []generated.Restaurant{
// 		{
// 			ID: 1, 
// 			Name: "Pizza Paradise", 
// 			Address: "123 Main Street", 
// 			Rating: pgtype.Numeric{Int: bigIntFromFloat(4.5), Valid: true}, 
// 			Category: ptr("pizza")},
// 		{
// 			ID: 2, 
// 			Name: "Burger Bonanza", 
// 			Address: "456 High Street", 
// 			Rating: pgtype.Numeric{Int: bigIntFromFloat(4.2), Valid: true},
// 			Category: ptr("burger")},
// 	}
// 	mockDomain.On("FetchAllRestaurants", mock.Anything).Return(expectedRestaurants, nil)

// 	// Create request and recorder
// 	req := httptest.NewRequest("GET", "/api/restaurants", nil)
// 	rec := httptest.NewRecorder()

// 	// Call handler
// 	handler.GetAllRestaurants()(rec, req)

// 	// Assertions
// 	mockDomain.AssertExpectations(t)
// 	if rec.Code != http.StatusOK {
// 		t.Errorf("expected status 200, got %d", rec.Code)
// 	}

// 	var actualRestaurants []generated.Restaurant
// 	err := json.Unmarshal(rec.Body.Bytes(), &actualRestaurants)
// 	if err != nil {
// 		t.Fatalf("failed to unmarshal response %v:", err)
// 		return
// 	}

// 	if len(actualRestaurants) != len(expectedRestaurants) {
// 		t.Errorf("expected %d restaurants, got %d", len(expectedRestaurants), len(actualRestaurants))
// 	}
// }

// func TestGetRestaurantById(t *testing.T) {
// 	mockDomain := new(MockRestaurantDomain)
// 	handler := NewRestaurantHandler(mockDomain)

// 	// Mock response
// 	expectedRestaurant := &generated.Restaurant{
// 		ID: 1,
// 		Name: "Pizza Paradise",
// 		Address: "123 Main Street",
// 		Rating: pgtype.Numeric{Int: bigIntFromFloat(4.5), Valid: true},
// 		Category: ptr("pizza"),
// 	}
// 	mockDomain.On("GetRestaurantById", mock.Anything, int32(1)).Return(expectedRestaurant, nil)

// 	// Create request and recorder
// 	req := httptest.NewRequest("GET", "/api/restaurants/1", nil)
// 	req = req.WithContext(context.WithValue(req.Context(), "restaurantId", int32(1))) // PathValue mock
// 	rec := httptest.NewRecorder()

// 	// Call handler
// 	handler.GetRestaurantById()(rec, req)

// 	// Assertions
// 	mockDomain.AssertExpectations(t)
// 	if rec.Code != http.StatusOK {
// 		t.Errorf("expected status 200, got %d", rec.Code)
// 	}

// 	var actualRestaurant generated.Restaurant
// 	err := json.Unmarshal(rec.Body.Bytes(), &actualRestaurant)
// 	if err != nil {
// 		t.Fatalf("failed to unmarshal response %v:", err)
// 		return
// 	}

// 	if actualRestaurant.ID != expectedRestaurant.ID {
// 		t.Errorf("expected restaurant ID %d, got %d", expectedRestaurant.ID, actualRestaurant.ID)
// 	}
// }

// func TestGetMenuItemsByRestaurant(t *testing.T) {
// 	mockDomain := new(MockRestaurantDomain)
// 	handler := NewRestaurantHandler(mockDomain)

// 	// Mock response
// 	expectedMenuitems := []generated.FetchMenuItemsByRestaurantIdRow{
// 		{
// 			ID: 1, 
// 			Name: "Pepperoni Pizza", 
// 			Description: ptr("Classic pizza with pepperoni."), 
// 			Price: pgtype.Numeric{Int: bigIntFromFloat(12.99), Valid: true}, 
// 			Restaurantid: 1},
// 		{
// 			ID: 2, 
// 			Name: "Margarita Pizza", 
// 			Description: ptr("Pizza with tomato, mozzarella, and basil."), 
// 			Price: pgtype.Numeric{Int: bigIntFromFloat(10.99), Valid: true}, 
// 			Restaurantid: 1},
// 	}
// 	mockDomain.On("FetchMenuItemsByRestaurantId", mock.MatchedBy(func(ctx context.Context) bool {return true}), int32(1)).Return(expectedMenuitems, nil)

// 	// Create request and recorder
// 	req := httptest.NewRequest("GET", "/api/restaurants/1/menu-items", nil)
// 	rec := httptest.NewRecorder()

// 	// Call handler
// 	handler.GetMenuItemsByRestaurant()(rec, req)

// 	// Assertions
// 	mockDomain.AssertExpectations(t)
// 	if rec.Code != http.StatusOK {
// 		t.Errorf("expected status 200, got %d", rec.Code)
// 	}

// 	var actualMenuitems []generated.FetchMenuItemsByRestaurantIdRow
// 	err := json.Unmarshal(rec.Body.Bytes(), &actualMenuitems)
// 	if err != nil {
// 		t.Fatalf("failed to unmarshal response: %v", err)
// 		return
// 	}

// 	if len(actualMenuitems) != len(expectedMenuitems) {
// 		t.Errorf("expected %d menuitems, got %d", len(expectedMenuitems), len(actualMenuitems))
// 	}
// }

// func TestFetchMenuItemByRestaurantAndId(t *testing.T) {
// 	mockDomain := new(MockRestaurantDomain)
// 	handler := NewRestaurantHandler(mockDomain)

// 	expectedMenuitem := &generated.FetchMenuItemsByRestaurantIdRow{
// 		ID: 1,
// 		Name: "Pepperoni Pizza",
// 		Description: ptr("Classic pizza with pepperoni"),
// 		Price: pgtype.Numeric{Int: bigIntFromFloat(12.99), Valid: true},
// 		Restaurantid: 1,
// 	}
// 	params := generated.GetMenuItemByRestaurantAndIdParams{
// 		Restaurantid: 1,
// 		ID: 1,
// 	}
// 	mockDomain.On("GetMenuItemByRestaurantAndId", mock.MatchedBy(func(ctx context.Context) bool {return true}), params).Return(expectedMenuitem, nil)

// 	req := httptest.NewRequest("GET", "/api/restaurants/1/menu-items/1", nil)
// 	req = req.WithContext(context.WithValue(req.Context(), "restaurantId", "1")) // PathValue mock
// 	req = req.WithContext(context.WithValue(req.Context(), "menuitemId", "1")) // PathValue mock
// 	rec := httptest.NewRecorder()

// 	handler.GetMenuItemByRestaurantAndId()(rec, req)

// 	mockDomain.AssertExpectations(t)
// 	if rec.Code != http.StatusOK {
// 		t.Errorf("expected status 200, got %d", rec.Code)
// 	}

// 	var actualMenuitem generated.FetchMenuItemsByRestaurantIdRow
// 	err := json.Unmarshal(rec.Body.Bytes(), &actualMenuitem)
// 	if err != nil {
// 		t.Fatalf("failed to unmarshal response: %v", err)
// 		return
// 	}

// 	if actualMenuitem.ID != expectedMenuitem.ID {
// 		t.Errorf("expected menuitem ID %d, got %d", expectedMenuitem.ID, actualMenuitem.ID)
// 	}
// }

// func TestFetchAllCategories(t *testing.T) {
// 	mockDomain := new(MockRestaurantDomain)
// 	handler := NewRestaurantHandler(mockDomain)

// 	expectedCategories := []string{"pizza", "burger", "sushi", "mexican"}
// 	mockDomain.On("FetchAllCategories", mock.Anything).Return(expectedCategories, nil)

// 	req := httptest.NewRequest("GET", "/api/categories", nil)
// 	rec := httptest.NewRecorder()

// 	handler.GetAllCategories()(rec, req)

// 	mockDomain.AssertExpectations(t)
// 	if rec.Code != http.StatusOK {
// 		t.Errorf("expected status 200, got %d", rec.Code)
// 	}

// 	var actualCategories []string
// 	err := json.Unmarshal(rec.Body.Bytes(), &actualCategories)
// 	if err != nil {
// 		t.Fatalf("failed to unmarshal response: %v", err)
// 		return
// 	}

// 	if len(actualCategories) != len(expectedCategories) {
// 		t.Errorf("expected %d categories, got %d", len(expectedCategories), len(actualCategories))
// 	}
// }

// func TestFilterRestaurantsByCategory(t *testing.T) {
// 	mockDomain := new(MockRestaurantDomain)
// 	handler := NewRestaurantHandler(mockDomain)

// 	expectedRestaurants := []generated.Restaurant{
// 		{
// 			ID: 1,
// 			Name: "Pizza Paradise",
// 			Address: "123 Main Street",
// 			Rating: pgtype.Numeric{Int: bigIntFromFloat(4.5), Valid: true},
// 			Category: ptr("pizza")},
// 	}
// 	mockDomain.On("FilterRestaurantsByCategory", mock.MatchedBy(func(ctx context.Context) bool {return true}), "pizza").Return(expectedRestaurants, nil)

// 	req := httptest.NewRequest("GET", "/api/filter/pizza", nil)
// 	rec := httptest.NewRecorder()

// 	handler.FilterRestaurantByCategory()(rec, req)

// 	mockDomain.AssertExpectations(t)
// 	if rec.Code != http.StatusOK {
// 		t.Errorf("expected status 200, got %d", rec.Code)
// 		return
// 	}

// 	var actualRestaurants []generated.Restaurant
// 	err := json.Unmarshal(rec.Body.Bytes(), &actualRestaurants)
// 	if err != nil {
// 		t.Fatalf("failed to unmarshal response: %v", err)
// 		return
// 	}

// 	if len(actualRestaurants) != len(expectedRestaurants) {
// 		t.Errorf("expected %d restaurants, got %d", len(expectedRestaurants), len(actualRestaurants))
// 	}
// }
