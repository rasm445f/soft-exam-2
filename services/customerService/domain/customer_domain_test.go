package domain_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/oTuff/go-startkode/db/generated"
	"github.com/oTuff/go-startkode/domain"
)

func TestCreateCustomer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error when opening sqlmock: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	queries := generated.New(db) // This will still fail without context support in sqlmock

	mock.ExpectExec("INSERT INTO customers").
		WithArgs("John Doe", "john.doe@example.com", "P@ssw0rd1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	customerDomain := domain.NewCustomerDomain(queries)

	// Test logic
	_, err = customerDomain.CreateCustomer(ctx, generated.CreateCustomerParams{
		Name:     pointer("John Doe"),
		Email:    pointer("john.doe@example.com"),
		Password: pointer("P@ssw0rd1"),
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
