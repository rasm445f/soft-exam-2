package domain

import (
	"context"
	"regexp"
	"testing"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/rasm445f/soft-exam-2/db/generated"
)

// Helper functions to create pointers for literals
func intPtr(i int32) *int32 {
	return &i
}
func stringPtr(s string) *string {
	return &s
}

func TestGetAllCustomersDomain(t *testing.T) {

	// Initialize pgxmock
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("failed to create pgxmock: %v", err)
	}
	defer mock.Close(context.Background())

	// Define the exact query to match
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

	// Define mock data with pointers where required
	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WillReturnRows(
			pgxmock.NewRows([]string{
				"id", "name", "email", "phonenumber", "password", "street_address", "zip_code", "city",
			}).
				AddRow(
					int32(1),
					stringPtr("John Doe"),             // Name as *string
					stringPtr("john.doe@example.com"), // Email as *string
					stringPtr("1234567890"),           // Phone number as *string
					stringPtr("hashedpassword"),       // Password as *string
					stringPtr("123 Main St"),          // Street address as *string
					intPtr(12345),                     // Zip code as *int32
					stringPtr("SomeCity"),             // City as *string
				).
				AddRow(
					int32(2),
					stringPtr("Jane Smith"),             // Name as *string
					stringPtr("jane.smith@example.com"), // Email as *string
					stringPtr("0987654321"),             // Phone number as *string
					stringPtr("hashedpassword"),         // Password as *string
					stringPtr("456 Elm St"),             // Street address as *string
					intPtr(67890),                       // Zip code as *int32
					stringPtr("OtherCity"),              // City as *string
				),
		)

	// Create sqlc Queries instance
	queries := generated.New(mock)

	// Create the domain instance
	customerDomain := NewCustomerDomain(queries)

	// Call the method under test
	customers, err := customerDomain.GetAllCustomersDomain(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Assert the results
	if len(customers) != 2 {
		t.Errorf("expected 2 customers, got %d", len(customers))
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestCreateCustomer(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("failed to create pgxmock: %v", err)
	}
	defer mock.Close(context.Background())

	queries := generated.New(mock)
	customerDomain := NewCustomerDomain(queries)

	// Test Case 1: Invalid Email
	t.Run("Invalid Email", func(t *testing.T) {
		params := generated.CreateCustomerParams{
			StreetAddress: stringPtr("123 Main St"),
			ZipCode:       intPtr(12345),
			Name:          stringPtr("John Doe"),
			Email:         stringPtr("invalid-email"),
			Phonenumber:   stringPtr("1234567890"),
			Password:      stringPtr("Password1!"),
		}

		err := customerDomain.CreateCustomerDomain(context.Background(), params)

		if err == nil || err.Error() != "invalid email format" {
			t.Errorf("expected error: invalid email format, got: %v", err)
		}
	})

	// Test Case 2: Invalid Password
	t.Run("Invalid Password", func(t *testing.T) {
		params := generated.CreateCustomerParams{
			StreetAddress: stringPtr("123 Main St"),
			ZipCode:       intPtr(12345),
			Name:          stringPtr("John Doe"),
			Email:         stringPtr("john.doe@example.com"),
			Phonenumber:   stringPtr("1234567890"),
			Password:      stringPtr("pass"),
		}

		err := customerDomain.CreateCustomerDomain(context.Background(), params)

		if err == nil || err.Error() != "password must be at least 8 characters long" {
			t.Errorf("expected error: password must be at least 8 characters long, got: %v", err)
		}
	})

	// Test Case 3: Valid Email and Password
	t.Run("Valid Email and Password", func(t *testing.T) {
		params := generated.CreateCustomerParams{
			StreetAddress: stringPtr("123 Main St"),
			ZipCode:       intPtr(12345),
			Name:          stringPtr("John Doe"),
			Email:         stringPtr("john.doe@example.com"),
			Phonenumber:   stringPtr("1234567890"),
			Password:      stringPtr("Password1!"),
		}

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

		err := customerDomain.CreateCustomerDomain(context.Background(), params)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
	})
}
