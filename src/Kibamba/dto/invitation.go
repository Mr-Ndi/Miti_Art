package dto

type InvitationInput struct {
	VendorEmail     string `json:"VendorEmail" example:"vendor@example.com" binding:"required"`
	VendorFirstName string `json:"VendorFirstName" example:"John" binding:"required"`
	VendorOtherName string `json:"VendorOtherName" example:"Doe" binding:"required"`
}

type InviteResponse struct {
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Invitation sent successfully"`
	SentTo  string `json:"sent_to" example:"vendor@example.com"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request data!"`
}
