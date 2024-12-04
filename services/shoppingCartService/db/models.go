package db

type ShoppingCart struct {
	CustomerId   int                `json:"customer_id"`
	RestaurantId int                `json:"restaurant_id"`
	TotalAmount  float64            `json:"total_amount"`
	VatAmount    float64            `json:"vat_amount"`
	Items        []ShoppingCartItem `json:"items"`
}

type ShoppingCartItem struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}
