package db

// TODO: implement the shoppingCart model as it is needed when creating a order
type ShoppingCart struct {
	Id           string             `json:"id"`
	UserId       int                `json:"user_id"`
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
