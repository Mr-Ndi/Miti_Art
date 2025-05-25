package dto

import "github.com/google/uuid"

type WishListRequest struct {
	ProductID uuid.UUID `json:"productID" binding:"required"`
}
