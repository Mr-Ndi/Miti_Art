package dto

type ProductResponse struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Material string  `json:"material"`
	Price    float64 `json:"price"`
	ImageURL string  `json:"imageUrl"`
}
