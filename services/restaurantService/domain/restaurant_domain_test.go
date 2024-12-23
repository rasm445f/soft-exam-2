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

	t.Run("Valid Restaurant Data", func(t *testing.T) {
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

	t.Run("No Restaurants", func(t *testing.T) {
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

func TestGetRestaurantByIdDomain(t *testing.T) {
	mock, _, domain := SetupTestMocks(t)
	defer CloseMocks(mock)

	t.Run("Valid ID", func(t *testing.T) {
		// Arrange
		row := pgxmock.NewRows([]string{"id", "name", "rating", "category", "address", "zip_code"}).
			AddRow(int32(1), "Pizza Paradise", float64Ptr(4.5), stringPtr("Pizza"), stringPtr("Main Street 123"), int32Ptr(2800))
		mock.ExpectQuery(`SELECT id, name, rating, category, address, zip_code FROM restaurant WHERE id = \$1`).
			WithArgs(int32(1)).
			WillReturnRows(row)

		// Act
		got, err := domain.GetRestaurantByIdDomain(context.Background(), int32(1))

		// Assert
		want := &generated.Restaurant{
			ID:       1,
			Name:     "Pizza Paradise",
			Rating:   float64Ptr(4.5),
			Category: stringPtr("Pizza"),
			Address:  stringPtr("Main Street 123"),
			ZipCode:  int32Ptr(2800),
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

	t.Run("Invalid Restaurant ID", func(t *testing.T) {
		got, err := domain.GetRestaurantByIdDomain(context.Background(), int32(0))

		if err == nil {
			t.Errorf("expected nil, got nil")
		}
		if got != nil {
			t.Errorf("expected nil, got: %+v", got)
		}
	})

	t.Run("Non-Existent ID", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, name, rating, category, address, zip_code FROM restaurant WHERE id = \$1`).
			WithArgs(int32(999)).
			WillReturnError(context.DeadlineExceeded)

		got, err := domain.GetRestaurantByIdDomain(context.Background(), int32(999))


		if err == nil || err.Error() != "restaurant not found" {
			t.Fatalf("expected error 'restaurant not found' got: %v", err)
		}
		if got != nil {
			t.Errorf("expected nil, got: %+v", got)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet mock expectations: %v", err)
		}
	})
}

func TestGetMenuItemsByRestaurantIdDomain(t *testing.T) {
	mock, _, domain := SetupTestMocks(t)
	defer CloseMocks(mock)

	t.Run("Valid Restaurant ID", func(t *testing.T) {
		// Arrange
		rows := pgxmock.NewRows([]string{"id", "restaurantid", "name", "price", "description"}).
			AddRow(int32(1), int32(1), "Cheese Pizza", float64(12.5), stringPtr("Delicious cheese pizza")).
			AddRow(int32(2), int32(1), "Veggie Pizza", float64(10.0), stringPtr("Healthy veggie pizza"))
		mock.ExpectQuery(`SELECT id, restaurantid, name, price, description FROM menuitem WHERE restaurantid = \$1`).
			WithArgs(int32(1)).
			WillReturnRows(rows)

		// Act
		got, err := domain.GetMenuItemsByRestaurantIdDomain(context.Background(), int32(1))

		// Assert
		want := []generated.Menuitem{
			{ID: 1, Restaurantid: 1, Name: "Cheese Pizza", Price: 12.5, Description: stringPtr("Delicious cheese pizza")},
			{ID: 2, Restaurantid: 1, Name: "Veggie Pizza", Price: 10.0, Description: stringPtr("Healthy veggie pizza")},
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

	t.Run("No MenuItems", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"id", "restaurantid", "name", "price", "description"})
		mock.ExpectQuery(`SELECT id, restaurantid, name, price, description FROM menuitem WHERE restaurantid = \$1`).
			WithArgs(int32(999)).
			WillReturnRows(rows)

		got, err := domain.GetMenuItemsByRestaurantIdDomain(context.Background(), int32(999))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("expected no menuitems, got %+v", got)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet mock expectations: %v", err)
		}
	})
	
	t.Run("Database Error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, restaurantid, name, price, description FROM menuitem WHERE restaurantid = \$1`).
			WithArgs(int32(1)).
			WillReturnError(context.DeadlineExceeded)

		got, err := domain.GetMenuItemsByRestaurantIdDomain(context.Background(), int32(1))

		if err == nil {
			t.Fatalf("expected an error, got nil")
		}
		if got != nil {
			t.Errorf("expected nil, got %+v", got)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet mock expectations: %v", err)
		}
	})
}

func TestGetMenuItemsByRestaurantAndIdDomain(t *testing.T) {
	mock, _, domain := SetupTestMocks(t)
	defer CloseMocks(mock)

	t.Run("Valid Restaurant ID and Menu Item ID", func(t *testing.T) {
		// Arrange
		rows := pgxmock.NewRows([]string{"id", "restaurantid", "name", "price", "description"}).
			AddRow(int32(1), int32(1), "Cheese Pizza", float64(12.5), stringPtr("Delicious cheese pizza"))
		mock.ExpectQuery(`SELECT id, restaurantid, name, price, description FROM menuitem WHERE restaurantid = \$1 AND id = \$2`).
			WithArgs(int32(1), int32(1)).
			WillReturnRows(rows)

		// Act
		params := generated.GetMenuItemByRestaurantAndIdParams{
			Restaurantid: int32(1),
			ID: int32(1),
		}
		got, err := domain.GetMenuItemByRestaurantAndIdDomain(context.Background(), params)

		// Assert
		want := &generated.Menuitem{
			ID:           1,
			Restaurantid: 1,
			Name:         "Cheese Pizza",
			Price:        12.5,
			Description:  stringPtr("Delicious cheese pizza"),
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

	t.Run("MenuItem Not Found", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, restaurantid, name, price, description FROM menuitem WHERE restaurantid = \$1 AND id = \$2`).
			WithArgs(int32(1), int32(99)).
			WillReturnError(context.DeadlineExceeded)

		params := generated.GetMenuItemByRestaurantAndIdParams{
			Restaurantid: int32(1),
			ID: int32(99),
		}
		got, err := domain.GetMenuItemByRestaurantAndIdDomain(context.Background(), params)

		if err == nil || err.Error() != "menuitem not found" {
			t.Fatalf("expected error 'menuitem not found', got: %v", err)
		}
		if got != nil {
			t.Errorf("expected nil, got: %+v", got)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet mock expectations: %v", err)
		}
	})
}

func TestGetAllCategoriesDomain(t *testing.T) {
	mock, _, domain := SetupTestMocks(t)
	defer CloseMocks(mock)
	
	t.Run("Valid Categories Data", func(t *testing.T) {
		// Arrange mock data
		rows := pgxmock.NewRows([]string{"category"}).
			AddRow(stringPtr("Pizza")).
			AddRow(stringPtr("Sushi"))
		mock.ExpectQuery(`SELECT DISTINCT category FROM restaurant WHERE category IS NOT NULL ORDER BY category`).
			WillReturnRows(rows)
		
		// Act
		got, err := domain.GetAllCategoriesDomain(context.Background())

		// Assert
		want := []string{"Pizza", "Sushi"}

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

	t.Run("No Categories", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"category"})
		mock.ExpectQuery(`SELECT DISTINCT category FROM restaurant WHERE category IS NOT NULL ORDER BY category`).
			WillReturnRows(rows)

		got, err := domain.GetAllCategoriesDomain(context.Background())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("expected no menuitems, got %+v", got)
		}
	})
}

func TestFilterRestaurantsByCategoryDomain(t *testing.T) {
	mock, _, domain := SetupTestMocks(t)
	defer CloseMocks(mock)

	t.Run("Valid Category", func(t *testing.T) {
		// Arrange
		rows := pgxmock.NewRows([]string{"id", "name", "rating", "category", "address", "zip_code"}).
			AddRow(int32(1), "Pizza Paradise", float64Ptr(4.5), stringPtr("Pizza"), stringPtr("Main Street 123"), int32Ptr(2800)).
			AddRow(int32(2), "Sushi World", float64Ptr(4.8), stringPtr("Pizza"), stringPtr("Second Street 456"), int32Ptr(2900))
		mock.ExpectQuery(`SELECT id, name, rating, category, address, zip_code FROM restaurant WHERE category ILIKE \$1`).
			WithArgs(stringPtr("Pizza")).
			WillReturnRows(rows)

		// Act
		got, err := domain.FilterRestaurantsByCategoryDomain(context.Background(), "Pizza")

		// Assert
		want := []generated.Restaurant{
			{ID: 1, Name: "Pizza Paradise", Rating: float64Ptr(4.5), Category: stringPtr("Pizza"), Address: stringPtr("Main Street 123"), ZipCode: int32Ptr(2800)},
			{ID: 2, Name: "Sushi World", Rating: float64Ptr(4.8), Category: stringPtr("Pizza"), Address: stringPtr("Second Street 456"), ZipCode: int32Ptr(2900)},
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

	t.Run("No Restaurants in Category", func(t *testing.T) {
		rows := pgxmock.NewRows([]string{"id", "name", "rating", "category", "address", "zip_code"})
		mock.ExpectQuery(`SELECT id, name, rating, category, address, zip_code FROM restaurant WHERE category ILIKE \$1`).
			WithArgs(stringPtr("NonExistent")).
			WillReturnRows(rows)

		got, err := domain.FilterRestaurantsByCategoryDomain(context.Background(), "NonExistent")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("expected no restaurants, got %+v", got)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet mock expectations: %v", err)
		}
	})
	
	t.Run("Database Error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, name, rating, category, address, zip_code FROM restaurant WHERE category ILIKE \$1`).
			WithArgs(stringPtr("NonExistent")).
			WillReturnError(context.DeadlineExceeded)

		got, err := domain.FilterRestaurantsByCategoryDomain(context.Background(), "Pizza")

		if err == nil {
			t.Fatalf("expected an error, got nil")
		}
		if got != nil {
			t.Fatalf("expected an nil, got %+v", got)
		}
	})
}

