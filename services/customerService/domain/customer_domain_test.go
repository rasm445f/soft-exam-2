package domain

import (
	"context"
	"regexp"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/rasm445f/soft-exam-2/db/generated"
)

// Assuming you have these helper functions and domain constructors
// in the same package, or imported from a helper file:
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

func float64Ptr(f float64) *float64 {
	return &f
}

func stringPtr(s string) *string {
	return &s
}

// ----------------------------------------------------------------------------
// CustomerDomain Tests
// ----------------------------------------------------------------------------

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
