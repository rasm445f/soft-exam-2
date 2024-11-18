package domain

import (
	"context"
	"errors"

	"github.com/oTuff/go-startkode/db"
)

type ShoppingCartDomain struct {
	repo *db.ShoppingCartRepository
}

func NewShoppingCartDomain(repo *db.ShoppingCartRepository) *ShoppingCartDomain {
	return &ShoppingCartDomain{repo: repo}
}

func (d *ShoppingCartDomain) AddItem(ctx context.Context, item db.AddItemParams) error {
	// Business validation
	if item.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	// Call repository
	return d.repo.AddItem(ctx, item)
}

func (d *ShoppingCartDomain) UpdateCart(ctx context.Context, userID string, itemID string, quantity int) error {
	// Business validation
	if quantity < 0 {
		return errors.New("quantity cannot be negative")
	}

	return d.repo.UpdateCart(ctx, userID, itemID, quantity)
}

func (d *ShoppingCartDomain) ViewCart(ctx context.Context, userID string) ([]db.ShoppingCartItem, error) {
	if len(userID) == 0 {
		return nil, errors.New("userID cannot be empty")
	}

	return d.repo.ViewCart(ctx, userID)
}
