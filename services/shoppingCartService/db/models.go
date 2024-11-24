package db

type ShoppingCart struct {
	CustomerId   int                `json:"customer_id"`
	RestaurantId int                `json:"restaurant_id"`
	TotalAmount  int                `json:"total_amount"`
	VatAmount    int                `json:"vat_amount"`
	Comment      string             `json:"comment"`
	Items        []ShoppingCartItem `json:"items"`
}

type ShoppingCartItem struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}
