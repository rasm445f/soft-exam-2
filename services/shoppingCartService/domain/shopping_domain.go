package domain

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/rasm445f/soft-exam-2/db"
)

type ShoppingCartDomain struct {
	repo *db.ShoppingCartRepository
}

func NewShoppingCartDomain(repo *db.ShoppingCartRepository) *ShoppingCartDomain {
	return &ShoppingCartDomain{repo: repo}
}

type AddItemParams struct {
	CustomerId   int     `json:"customerId"`
	RestaurantId int     `json:"restaurantId"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Quantity     int     `json:"quantity"`
}

// Helper function to calculate cart totals
func (d *ShoppingCartDomain) recalculateCartTotals(cart *db.ShoppingCart) {
	cart.TotalAmount = 0
	for _, item := range cart.Items {
		cart.TotalAmount += item.Price * float64(item.Quantity)
	}
	cart.VatAmount = cart.TotalAmount * 0.20
}

func (d *ShoppingCartDomain) AddItemDomain(ctx context.Context, itemParams AddItemParams) error {
	// Business validation
	if itemParams.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	// check if cart already exist
	cart, err := d.repo.GetCart(ctx, itemParams.CustomerId)
	if err == redis.Nil {
		// if it does not exist create the cart
		cart = &db.ShoppingCart{
			CustomerId:   itemParams.CustomerId,
			RestaurantId: itemParams.RestaurantId,
			Items:        []db.ShoppingCartItem{},
		}
	} else if err != nil {
		return err
	}

	// Add item
	item := db.ShoppingCartItem{
		Id:       len(cart.Items) + 1, // Simple ID generation instead of redis INCR command
		Name:     itemParams.Name,
		Price:    itemParams.Price,
		Quantity: itemParams.Quantity,
	}

	cart.Items = append(cart.Items, item)

	d.recalculateCartTotals(cart)

	return d.repo.SaveCart(ctx, cart)
}

func (d *ShoppingCartDomain) UpdateCartDomain(ctx context.Context, customerId, itemID, quantity int) error {
	if quantity < 0 {
		return errors.New("quantity cannot be negative")
	}

	cart, err := d.repo.GetCart(ctx, customerId)
	if err != nil {
		return err
	}

	// Update item quantity
	found := false
	for i := range cart.Items {
		if cart.Items[i].Id == itemID {
			if quantity == 0 {
				// Remove item
				cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			} else {
				cart.Items[i].Quantity = quantity
			}
			found = true
			break
		}
	}
	if !found {
		return errors.New("item not found in cart")
	}

	d.recalculateCartTotals(cart)
	return d.repo.SaveCart(ctx, cart)
}

func (d *ShoppingCartDomain) ViewCartDomain(ctx context.Context, customerId int) (*db.ShoppingCart, error) {
	return d.repo.GetCart(ctx, customerId)
}

func (d *ShoppingCartDomain) ClearCartDomain(ctx context.Context, customerId int) error {
	return d.repo.ClearCart(ctx, customerId)
}
