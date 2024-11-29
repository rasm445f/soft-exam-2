package domain

import (
	"context"
	"errors"

	"github.com/rasm445f/soft-exam-2/db"
)

type ShoppingCartDomain struct {
	repo *db.ShoppingCartRepository
}

func NewShoppingCartDomain(repo *db.ShoppingCartRepository) *ShoppingCartDomain {
	return &ShoppingCartDomain{repo: repo}
}

func (d *ShoppingCartDomain) AddItemDomain(ctx context.Context, item db.AddItemParams) error {
	// Business validation
	if item.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}

	// Call repository
	return d.repo.AddItem(ctx, item)
}

func (d *ShoppingCartDomain) UpdateCartDomain(ctx context.Context, customerId int, itemID int, quantity int) error {
	// Business validation
	if quantity < 0 {
		return errors.New("quantity cannot be negative")
	}

	return d.repo.UpdateCart(ctx, customerId, itemID, quantity)
}

func (d *ShoppingCartDomain) ViewCartDomain(ctx context.Context, costumerId int) (*db.ShoppingCart, error) {
	// if costumerId == nil {
	// 	return nil, errors.New("customerId cannot be empty")
	// }

	return d.repo.ViewCart(ctx, costumerId)
}
func (d *ShoppingCartDomain) ClearCartDomain(ctx context.Context, costumerId int) error {
	// if costumerId == nil {
	// 	return nil, errors.New("customerId cannot be empty")
	// }

	return d.repo.ClearCart(ctx, costumerId)
}
