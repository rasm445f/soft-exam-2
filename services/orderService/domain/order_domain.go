package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/rasm445f/soft-exam-2/db/generated"
)

type OrderDomain struct {
	repo *generated.Queries
}

// NewOrderDomain initializes the domain layer
func NewOrderDomain(repo *generated.Queries) *OrderDomain {
	return &OrderDomain{repo: repo}
}

func (d *OrderDomain) GetAllOrdersDomain(ctx context.Context) ([]generated.Order, error) {
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

func (d *OrderDomain) GetOrderByIdDomain(ctx context.Context, orderId int32) (*generated.Order, error) {
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

func (d *OrderDomain) CreateOrderDomain(ctx context.Context, orderParams generated.CreateOrderParams) (int32, error) {	
	amountExcludingVAT := orderParams.Totalamount - orderParams.Vatamount

	feeid, err := d.CalculateFee(ctx, amountExcludingVAT)
	if err != nil {
		fmt.Printf("%v", err)
		return 0, err
	}
	orderParams.Feeid = &feeid

	// Call the repository layer to create the order
	orderid, err := d.repo.CreateOrder(ctx, orderParams)
	if err != nil {
		return 0, errors.New("failed to create order: " + err.Error())
	}

	return orderid, nil
}

func (d *OrderDomain) UpdateOrderStatusDomain(ctx context.Context, orderId int32, status string) error {
	// Call the repository layer to update the order
	err := d.repo.UpdateOrderStatus(ctx, generated.UpdateOrderStatusParams{
		Status: status,
		ID: orderId,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("order not found")
		}
		return err
	}

	return nil
}

func (d *OrderDomain) DeleteOrderDomain(ctx context.Context, orderId int32) (error) {
	err := d.repo.DeleteOrderItemsByOrderId(ctx, orderId)
	if err != nil {
		return errors.New("failed to delete order items: " + err.Error())
	}

	err = d.repo.DeleteOrder(ctx, orderId)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("order not found")
		}
		return errors.New("failed to delete order: " + err.Error())
	}
	return nil
}

func (d *OrderDomain) CreateOrderItemDomain(ctx context.Context, itemParams generated.CreateOrderItemParams) (int32, error) {
	// Call the repository layer to create the order
	itemid, err := d.repo.CreateOrderItem(ctx, itemParams)
	if err != nil {
		return 0, errors.New("failed to create order item: " + err.Error())
	}

	return itemid, nil
}

func (d *OrderDomain) CalculateFee(ctx context.Context, amount float64) (int32, error) {
	var fee float64
	var percent float64

	if amount <= 100 {
		percent = 0.06
		fee = amount * percent
	} else if amount > 100 && amount <= 500 {
		percent = 0.05
		fee = amount * percent
	} else if amount > 500 && amount <= 1000 {
		percent = 0.04
		fee = amount * percent
	} else {
		percent = 0.03
		fee = amount * percent
	}

	desc := "some description"
	newFee := generated.CreateFeeParams{
		Percentage: &percent,
		Amount: &fee,
		Description: &desc ,
	}

	feeid, err := d.repo.CreateFee(ctx, newFee)
	if err != nil {
		return 0, errors.New("failed to create fee: " + err.Error())
	}

	return feeid, nil
}


// TODO: implement
// create fee
// create bonus
// create payment