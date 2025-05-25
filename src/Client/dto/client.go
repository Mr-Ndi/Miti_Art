package dto

import (
	"github.com/google/uuid"
)

// ClientRegisterRequest represents the payload for client registration.
type ClientRegisterRequest struct {
	ClientFirstName string `json:"clientFirstName" binding:"required"`
	ClientOtherName string `json:"clientOtherName" binding:"required"`
	ClientEmail     string `json:"clientEmail" binding:"required,email"`
	ClientPassword  string `json:"clientPassword" binding:"required"`
}

// LoginRequest represents a generic login request.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// OrderRequest represents the payload for creating an order.
type OrderRequest struct {
	ProductID uuid.UUID `json:"productID" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required"`
}

// WishlistRequest represents the payload for adding a product to the wishlist.
type WishlistRequest struct {
	ProductID uuid.UUID `json:"productID" binding:"required"`
}
