package dto

// LoginRequest represents login input payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required,min=6" example:"secret123"`
}

// LoginResponse represents login output (token)
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOi..."`
}
