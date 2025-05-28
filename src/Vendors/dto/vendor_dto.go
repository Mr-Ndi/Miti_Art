package dto

type VendorRegisterRequest struct {
	VendorPassword string `json:"vendorPassword" binding:"required"`
	VendorTin      int    `json:"vendorTin" binding:"required"`
	ShopName       string `json:"ShopName" binding:"required"`
}

type RegisterResponse struct {
	Message     string `json:"message"`
	VendorEmail string `json:"vendorEmail"`
}

// type ErrorResponse struct {
// 	Error string `json:"error"`
// }

type EditProductRequest struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
	Material string  `json:"material"`
	ImageURL string  `json:"image_url"`
}
