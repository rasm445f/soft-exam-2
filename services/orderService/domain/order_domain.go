package domain

import (
	"context"
	"database/sql"
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

func (d *OrderDomain) FetchAllOrders(ctx context.Context) ([]generated.Order, error) {
	rows, err := d.repo.GetAllOrders(ctx)
	if err != nil {
		return nil, errors.New("failed to fetch orders")
	}

	var orders []generated.Order
	for _, row := range rows {
		orders = append(orders, generated.Order{
			ID: row.ID,
			Totalamount: row.Totalamount,
			Vatamount: row.Vatamount,
			Status: row.Status,
			Timestamp: row.Timestamp,
			Comment: row.Comment,
			Customerid: row.Customerid,
			Restaurantid: row.Restaurantid,
			Deliveryagentid: row.Deliveryagentid,
			Paymentid: row.Paymentid,
			Bonusid: row.Bonusid,
			Feeid: row.Feeid,
		})
	}
	return orders, nil
}

func (d *OrderDomain) GetOrderById(ctx context.Context, orderId int32) (*generated.Order, error) {
	if orderId <= 0 {
		return nil, errors.New("invalid order id")
	}

	row, err := d.repo.GetOrderById(ctx, orderId)
	if err != nil {
		return nil, errors.New("order not found")
	}

	order := &generated.Order{
		ID: row.ID,
		Totalamount: row.Totalamount,
		Vatamount: row.Vatamount,
		Status: row.Status,
		Timestamp: row.Timestamp,
		Comment: row.Comment,
		Customerid: row.Customerid,
		Restaurantid: row.Restaurantid,
		Deliveryagentid: row.Deliveryagentid,
		Paymentid: row.Paymentid,
		Bonusid: row.Bonusid,
		Feeid: row.Feeid,
	}

	return order, nil
}

func (o *OrderDomain) CreateOrder(ctx context.Context, orderParams generated.CreateOrderParams) (int32, error) {
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

func (o *OrderDomain) DeleteOrder(ctx context.Context, orderId int32) (error) {
	err := o.repo.DeleteOrderItemsByOrderId(ctx, orderId)
	if err != nil {
		return errors.New("failed to delete order items: " + err.Error())
	}

	err = o.repo.DeleteOrder(ctx, orderId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("order not found")
		}
		return errors.New("failed to delete order: " + err.Error())
	}
	return nil
}


// TODO: implement
// create fee
// create bonus
// create payment