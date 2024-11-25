package broker

// Event represents the structure of a message to be published/consumed.
type Event struct {
	Type string `json:"type"`			// Event type, e.g., "order.placed"
	Payload interface{} `json:"payload"`	// Event payload, e.g., order details
}

// Event Types
const (
	MenuItemSelected = "menu_item_selected"
	CartUpdated = "cart_updated"
	OrderCreated = "order_created"
)