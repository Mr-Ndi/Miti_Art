package dto

type InvitationInput struct {
	VendorEmail     string `json:"VendorEmail" example:"vendor@example.com" binding:"required"`
	VendorFirstName string `json:"VendorFirstName" example:"Ninshuti" binding:"required"`
	VendorOtherName string `json:"VendorOtherName" example:"Poli" binding:"required"`
}
