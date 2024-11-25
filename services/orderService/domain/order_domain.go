package domain

import (
	"context"
	"errors"

	"github.com/rasm445f/soft-exam-2/db/generated"
)

type OrderDomain struct {
	repo *generated.Queries
}

// NewOrderDomain initializes the domain layer
func NewOrderDomain(repo *generated.Queries) *OrderDomain {
	return &OrderDomain{repo: repo}
}

func (o *OrderDomain) CreateOrder(ctx context.Context, orderParams generated.CreateOrderParams) (int32, error) {
	// Validate required fields
	// if orderParams.TotalAmount <= 0 {
	// 	return 0, errors.New("total amount must be greater than zero")
	// }
	// if orderParams.Status == "" {
	// 	return 0, errors.New("order status cannot be empty")
	// }
	
	// Call the repository layer to create the order
	orderid, err := o.repo.CreateOrder(ctx, orderParams)
	if err != nil {
		return 0, errors.New("failed to create order: " + err.Error())
	}

	return orderid, nil
}

func (o *OrderDomain) CreateOrderItem(ctx context.Context, itemParams generated.CreateOrderItemParams) (int32, error) {
		// Call the repository layer to create the order
		itemid, err := o.repo.CreateOrderItem(ctx, itemParams)
		if err != nil {
			return 0, errors.New("failed to create order item: " + err.Error())
		}
	
		return itemid, nil
}



// TODO: implement
// create fee
// create bonus
// create payment