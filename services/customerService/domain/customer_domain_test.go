package domain

import (
	"context"
	"database/sql"
	"regexp"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/rasm445f/soft-exam-2/db/generated"
)

func SetupTestMocks(t *testing.T) (pgxmock.PgxPoolIface, *generated.Queries, *CustomerDomain) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create pgxmock pool: %v", err)
	}

	queries := generated.New(mock)
	domain := NewCustomerDomain(queries)

	return mock, queries, domain
}

func CloseMocks(mock pgxmock.PgxPoolIface) {
	mock.Close()
}

func int32Ptr(i int32) *int32 {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func TestGetAllCustomersDomain(t *testing.T) {
	mock, _, customerDomain := SetupTestMocks(t)
	defer CloseMocks(mock)

	t.Run("Valid Data", func(t *testing.T) {
		// Arrange
		// Use your actual query string here
		expectedQuery := `SELECT 
            c.id,
            c.name,
            c.email,
            c.phonenumber,
            c.password,
            a.street_address AS street_address,
            z.zip_code,
            z.city
        FROM customer c
        LEFT JOIN address a ON c.addressid = a.id
        LEFT JOIN zipcode z ON a.zip_code = z.zip_code
        ORDER BY c.name`

		rows := pgxmock.NewRows([]string{
			"id", "name", "email", "phonenumber", "password", "street_address", "zip_code", "city",
		}).
			AddRow(
				int32(1),
				stringPtr("Alice Wonderland"),
				stringPtr("alice@example.com"),
				stringPtr("1234567890"),
				stringPtr("somehashedpass"),
				stringPtr("123 Main St"),
				int32Ptr(12345),
				stringPtr("ExampleCity"),
			).
			AddRow(
				int32(2),
				stringPtr("Bob Builder"),
				stringPtr("bob@example.com"),
				stringPtr("0987654321"),
				stringPtr("anotherhashedpass"),
				stringPtr("456 Elm St"),
				int32Ptr(67890),
				stringPtr("OtherCity"),
			)

		// Expect query and return rows
		mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnRows(rows)

		// Act
		got, err := customerDomain.GetAllCustomersDomain(context.Background())

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 2 {
			t.Errorf("expected 2 customers, got %d", len(got))
		}

		// (Optional) Example: check specific fields
		// want := []generated.Customer{ ... }
		// if !reflect.DeepEqual(got, want) { ... }

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet mock expectations: %v", err)
		}
	})

	t.Run("No Data", func(t *testing.T) {
		// Arrange
		expectedQuery := `SELECT 
            c.id,
            c.name,
            c.email,
            c.phonenumber,
            c.password,
            a.street_address AS street_address,
            z.zip_code,
            z.city
        FROM customer c
        LEFT JOIN address a ON c.addressid = a.id
        LEFT JOIN zipcode z ON a.zip_code = z.zip_code
        ORDER BY c.name`

		// Return an empty result set
		emptyRows := pgxmock.NewRows([]string{
			"id", "name", "email", "phonenumber", "password", "street_address", "zip_code", "city",
		})

		mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnRows(emptyRows)

		// Act
		got, err := customerDomain.GetAllCustomersDomain(context.Background())

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != 0 {
			t.Errorf("expected 0 customers, got %d", len(got))
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet mock expectations: %v", err)
		}
	})

	t.Run("Database Error", func(t *testing.T) {
		// Arrange
		expectedQuery := `SELECT 
            c.id,
            c.name,
            c.email,
            c.phonenumber,
            c.password,
            a.street_address AS street_address,
            z.zip_code,
            z.city
        FROM customer c
        LEFT JOIN address a ON c.addressid = a.id
        LEFT JOIN zipcode z ON a.zip_code = z.zip_code
        ORDER BY c.name`

		mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WillReturnError(context.DeadlineExceeded)

		// Act
		got, err := customerDomain.GetAllCustomersDomain(context.Background())

		// Assert
		if err == nil {
			t.Fatalf("expected an error but got nil")
		}
		if got != nil {
			t.Errorf("expected nil result, got %+v", got)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet mock expectations: %v", err)
		}
	})
}

func TestCreateCustomerDomain(t *testing.T) {
	mock, _, customerDomain := SetupTestMocks(t)
	defer CloseMocks(mock)

	t.Run("Valid Email & Password", func(t *testing.T) {
		// Arrange
		params := generated.CreateCustomerParams{
			StreetAddress: stringPtr("123 Main St"),
			ZipCode:       int32Ptr(12345),
			Name:          stringPtr("Charlie Chaplin"),
			Email:         stringPtr("charlie@example.com"),
			Phonenumber:   stringPtr("1234567890"),
			Password:      stringPtr("Password1!"),
		}

		// Expect the INSERT with proper arguments
		mock.ExpectExec(regexp.QuoteMeta(`WITH new_address AS (
            INSERT INTO address (street_address, zip_code)
            VALUES ($1, $2)
            RETURNING id AS address_id
        )
        INSERT INTO customer (name, email, phonenumber, addressid, password)
        VALUES ($3, $4, $5, (SELECT address_id FROM new_address), $6)`)).
			WithArgs(
				params.StreetAddress,
				params.ZipCode,
				params.Name,
				params.Email,
				params.Phonenumber,
				params.Password,
			).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		// Act
		err := customerDomain.CreateCustomerDomain(context.Background(), params)

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet mock expectations: %v", err)
		}
	})

	t.Run("Invalid Email", func(t *testing.T) {
		// Arrange
		params := generated.CreateCustomerParams{
			StreetAddress: stringPtr("123 Main St"),
			ZipCode:       int32Ptr(12345),
			Name:          stringPtr("Invalid Email Person"),
			Email:         stringPtr("not-an-email"),
			Phonenumber:   stringPtr("1234567890"),
			Password:      stringPtr("Password1!"),
		}

		// Act
		err := customerDomain.CreateCustomerDomain(context.Background(), params)

		// Assert
		if err == nil {
			t.Fatalf("expected invalid email error, got nil")
		}
		// You can check the exact error message if you throw custom errors
		if err.Error() != "invalid email format" {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("Invalid Password", func(t *testing.T) {
		// Arrange
		params := generated.CreateCustomerParams{
			StreetAddress: stringPtr("123 Main St"),
			ZipCode:       int32Ptr(12345),
			Name:          stringPtr("Weak Password Person"),
			Email:         stringPtr("weak.pass@example.com"),
			Phonenumber:   stringPtr("1234567890"),
			Password:      stringPtr("123"), // obviously too short
		}

		// Act
		err := customerDomain.CreateCustomerDomain(context.Background(), params)

		// Assert
		if err == nil {
			t.Fatalf("expected invalid password error, got nil")
		}
		if err.Error() != "password must be at least 8 characters long" {
			t.Errorf("unexpected error: %v", err)
		}
	})
}

func TestGetCustomerByIdDomain(t *testing.T) {
	mock, _, domain := SetupTestMocks(t)
	defer CloseMocks(mock)

	t.Run("Valid ID", func(t *testing.T) {
		// Arrange
		expectedQuery := `SELECT 
    c.id,
    c.name,
    c.email,
    c.phonenumber,
    c.password,
    a.street_address,
    z.zip_code,
    z.city
FROM customer c
LEFT JOIN address a ON c.addressid = a.id
LEFT JOIN zipcode z ON a.zip_code = z.zip_code
WHERE c.id = $1`
		id := int32(1)

		// Mock returning a valid customer with address and zip code
		row := pgxmock.NewRows([]string{
			"id", "name", "email", "phonenumber", "password", "street_address", "zip_code", "city"}).
			AddRow(
				int32(1),
				stringPtr("Alice Wonderland"),
				stringPtr("alice@example.com"),
				stringPtr("1234567890"),
				stringPtr("hashedpassword"),
				stringPtr("123 Main St"),
				int32Ptr(12345),
				stringPtr("Wonderland City"),
			)

		mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(id).WillReturnRows(row)

		// Act
		customer, err := domain.GetCustomerByIdDomain(context.Background(), id)

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if customer.ID != id {
			t.Errorf("expected ID %d, got %d", id, customer.ID)
		}
		if *customer.Name != "Alice Wonderland" {
			t.Errorf("expected name %s, got %s", "Alice Wonderland", *customer.Name)
		}
		if *customer.Email != "alice@example.com" {
			t.Errorf("expected email %s, got %s", "alice@example.com", *customer.Email)
		}
		if *customer.Phonenumber != "1234567890" {
			t.Errorf("expected phone number %s, got %s", "1234567890", *customer.Phonenumber)
		}
		if *customer.Password != "hashedpassword" {
			t.Errorf("expected password %s, got %s", "hashedpassword", *customer.Password)
		}
		if *customer.StreetAddress != "123 Main St" {
			t.Errorf("expected street address %s, got %s", "123 Main St", *customer.StreetAddress)
		}
		if *customer.ZipCode != 12345 {
			t.Errorf("expected zip code %d, got %d", 12345, *customer.ZipCode)
		}
		if *customer.City != "Wonderland City" {
			t.Errorf("expected city %s, got %s", "Wonderland City", *customer.City)
		}
	})
}

func TestDeleteCustomerDomain(t *testing.T) {
	mock, _, customerDomain := SetupTestMocks(t)
	defer CloseMocks(mock)

	t.Run("Valid Deletion", func(t *testing.T) {
		// Arrange
		expectedQuery := `DELETE FROM customer WHERE id = $1`
		id := int32(1)

		mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).WithArgs(id).WillReturnResult(pgxmock.NewResult("DELETE", 1))

		// Act
		err := customerDomain.DeleteCustomerDomain(context.Background(), id)

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet mock expectations: %v", err)
		}
	})

	t.Run("Non-existent ID", func(t *testing.T) {
		// Arrange
		expectedQuery := `DELETE FROM customer WHERE id = $1`
		id := int32(999)

		mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).WithArgs(id).WillReturnError(sql.ErrNoRows)

		// Act
		err := customerDomain.DeleteCustomerDomain(context.Background(), id)

		// Assert
		if err == nil || err.Error() != "sql: no rows in result set" {
			t.Fatalf("expected sql.ErrNoRows, got %v", err)
		}
	})
}

func TestUpdateCustomerDomain(t *testing.T) {
	mock, _, customerDomain := SetupTestMocks(t)
	defer CloseMocks(mock)

	t.Run("Valid Update", func(t *testing.T) {
		// Arrange
		params := generated.UpdateCustomerParams{
			ID:          int32(1),
			Name:        stringPtr("Updated Name"),
			Email:       stringPtr("updated@example.com"),
			Phonenumber: stringPtr(""),
			Password:    stringPtr("Updated1!"),
		}

		expectedQuery := `UPDATE customer
SET 
    name = COALESCE($2, name),
    email = COALESCE($3, email),
    phonenumber = COALESCE($4, phonenumber),
    password = COALESCE($5, password)
WHERE id = $1`

		mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
			WithArgs(
				params.ID,
				params.Name,
				params.Email,
				params.Phonenumber,
				params.Password,
			).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		// Act
		err := customerDomain.UpdateCustomerDomain(context.Background(), params)

		// Assert
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet mock expectations: %v", err)
		}
	})

	// 	t.Run("Non-existent ID", func(t *testing.T) {
	// 		// Arrange
	// 		params := generated.UpdateCustomerParams{
	// 			ID:          int32(999),
	// 			Name:        stringPtr("Updated Name"),
	// 			Email:       stringPtr("nonexistent@example.com"),
	// 			Password:    stringPtr("Updated1!"),
	// 			Phonenumber: stringPtr("1234567890"),
	// 		}

	// 		expectedQuery := `UPDATE customer
	// SET
	//     name = COALESCE($2, name),
	//     email = COALESCE($3, email),
	//     phonenumber = COALESCE($4, phonenumber),
	//     password = COALESCE($5, password)
	// WHERE id = $1`

	// 		mock.ExpectExec(regexp.QuoteMeta(expectedQuery)).
	// 			WithArgs(
	// 				params.ID,
	// 				params.Name,
	// 				params.Email,
	// 				params.Phonenumber,
	// 				params.Password,
	// 			).
	// 			WillReturnResult(pgxmock.NewResult("UPDATE", 0)) // No rows updated

	// 		// Act
	// 		err := customerDomain.UpdateCustomerDomain(context.Background(), params)

	//		// Assert
	//		if err == nil || err.Error() != "customer not found" {
	//			t.Fatalf("expected 'customer not found', got %v", err)
	//		}
	//	})
}
