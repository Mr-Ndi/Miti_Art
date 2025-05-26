package dto

type OrderResponse struct {
	OrderID     string  `json:"orderId"`
	ProductID   string  `json:"productId"`
	ProductName string  `json:"productName"`
	Quantity    int     `json:"quantity"`
	TotalPrice  float64 `json:"totalPrice"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"createdAt"`
}
