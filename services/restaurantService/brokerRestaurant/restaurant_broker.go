package brokerrestaurant

import (
	"context"
	"log"

	"github.com/rasm445f/soft-exam-2/broker"
	"github.com/rasm445f/soft-exam-2/db/generated"
	"github.com/rasm445f/soft-exam-2/domain"
)

type SelectItemParams struct {
	CustomerId   int32 `json:"customerId" example:"1"`
	RestaurantId int32 `json:"restaurantId" example:"10"`
	ItemId       int32 `json:"id"`
	Quantity     int   `json:"quantity" example:"2"`
}

type MenuItemSelection struct {
	CustomerID   int32   `json:"customerId" example:"1"`
	RestaurantId int32   `json:"restaurantId" example:"10"`
	Name         string  `json:"name" example:"Cheese Burger"`
	Price        float64 `json:"price" example:"10.00"`
	Quantity     int     `json:"quantity" example:"2"`
}

func SelectMenuBroker(selectionParams SelectItemParams, ctx context.Context, domain *domain.RestaurantDomain) error {
	var menuSelectionParams = generated.GetMenuItemByRestaurantAndIdParams{
		Restaurantid: selectionParams.RestaurantId,
		ID:           selectionParams.ItemId,
	}

	log.Printf("Fetching menu item with params: %+v", menuSelectionParams)
	intermediateMenuItem, err := domain.GetMenuItemByRestaurantAndIdDomain(ctx, menuSelectionParams)
	if err != nil {
		log.Printf("Error fetching menu item: %v", err)
		return err
	}

	menuItemSelection := MenuItemSelection{
		CustomerID:   selectionParams.CustomerId,
		RestaurantId: intermediateMenuItem.Restaurantid,
		Name:         intermediateMenuItem.Name,
		Price:        intermediateMenuItem.Price,
		Quantity:     selectionParams.Quantity,
	}

	event := broker.Event{
		Type:    broker.MenuItemSelected,
		Payload: menuItemSelection,
	}

	log.Printf("Publishing event to RabbitMQ: %+v", event)
	err = broker.Publish("menu_item_selected_queue", event)
	if err != nil {
		log.Printf("RabbitMQ publish failed: %v", err)
		return err
	}
	return nil
}
