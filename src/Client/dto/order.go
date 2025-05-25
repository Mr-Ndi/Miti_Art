// src/Client/dto/order.go
package dto

import "github.com/google/uuid"

type CreateOrderRequest struct {
	ProductID uuid.UUID `json:"productID" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required"`
}
