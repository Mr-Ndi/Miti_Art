package dto

type ProductResponse struct {
	ID       string  `json:"id"` // converted from uuid.UUID
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
	Material string  `json:"material"`
	ImageURL string  `json:"image_url"`
	VendorID string  `json:"vendor_id"` // also converted from uuid.UUID
}
