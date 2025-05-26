package dto

type ProductResponse struct {
	ID       string  `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name     string  `json:"name" example:"Wooden Sculpture"`
	Category string  `json:"category" example:"Art"`
	Material string  `json:"material" example:"Oak Wood"`
	Price    float64 `json:"price" example:"150.75"`
	ImageURL string  `json:"imageUrl" example:"https://res.cloudinary.com/demo/image/upload/sample.jpg"`
}
