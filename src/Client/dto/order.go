// src/Client/dto/order.go
package dto

import "github.com/google/uuid"

type CreateOrderRequest struct {
	ProductID uuid.UUID `json:"productID" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required"`
}
type OrderQuantityRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}
type ErrorResponse struct {
	Error string `json:"error"`
}
