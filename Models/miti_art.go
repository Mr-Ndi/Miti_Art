package models

import (
	"github.com/google/uuid"
)

// Role Enum
type Role string

const (
	RoleAdmin    Role = "ADMIN"
	RoleVendor   Role = "VENDOR"
	RoleCustomer Role = "CUSTOMER"
)

// Status Enum
type Status string

const (
	StatusPending   Status = "PENDING"
	StatusShipped   Status = "SHIPPED"
	StatusCompleted Status = "COMPLETED"
)

// User Model
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FirstName string    `gorm:"not null"`
	OtherName string    `gorm:"not null"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Salt      string    `gorm:"not null"`
	Role      Role      `gorm:"type:varchar(20);default:'CUSTOMER'"`

	// Relations
	Vendor   *Vendor    `gorm:"foreignKey:UserID"`
	Orders   []Order    `gorm:"foreignKey:UserID"`
	Wishlist []Wishlist `gorm:"foreignKey:UserID"`
}

// Vendor Model
type Vendor struct {
	UserID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	User         User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	BusinessName string    `gorm:"not null"`
	TaxPin       int64     `gorm:"unique;not null"`
	Approved     bool      `gorm:"default:false"`

	// Relations
	Products []Product `gorm:"foreignKey:VendorID"`
}

// Product Model
type Product struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	VendorID uuid.UUID `gorm:"not null"`
	Vendor   Vendor    `gorm:"foreignKey:VendorID;constraint:OnDelete:CASCADE"`
	Name     string    `gorm:"not null"`
	Price    float64   `gorm:"not null"`
	Category string    `gorm:"not null"`
	ImageURL string    `gorm:"not null"`

	// Relations
	Orders   []Order    `gorm:"foreignKey:ProductID"`
	Wishlist []Wishlist `gorm:"foreignKey:ProductID"`
}

// Order Model
type Order struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"not null"`
	User      User      `gorm:"references:ID;constraint:OnDelete:CASCADE"`
	ProductID uuid.UUID `gorm:"not null"`
	Product   Product   `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
	Quantity  int       `gorm:"not null"`
	Status    Status    `gorm:"type:varchar(20);default:'PENDING'"`
}

// Wishlist Model
type Wishlist struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	ProductID uuid.UUID `gorm:"not null"`
	Product   Product   `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}
