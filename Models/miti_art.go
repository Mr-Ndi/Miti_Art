package models

import (
	"time"

	"github.com/google/uuid"
)

// For payment
type PaymentStatus string

const (
	PaymentPending PaymentStatus = "PENDING"
	PaymentSuccess PaymentStatus = "SUCCESS"
	PaymentFailed  PaymentStatus = "FAILED"
)

// PaymentMethod Enum
type PaymentMethod string

const (
	MethodMomo   PaymentMethod = "MOMO"
	MethodCard   PaymentMethod = "CARD"
	MethodBank   PaymentMethod = "BANK"
	MethodPayPal PaymentMethod = "PAYPAL"
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

// Category Enum
type Category string

const (
	CategoryKitchen    Category = "Kitchen"
	CategoryLivingroom Category = "Living Room"
	CategoryBedroom    Category = "Bedroom"
	CategoryOffice     Category = "Office"
	CategoryOutdoor    Category = "Outdoor"
)

// Material Enum
type Material string

const (
	MaterialEucalyptus      Material = "Eucalyptus"
	MaterialGrevillea       Material = "Grevillea"
	MaterialPine            Material = "Pine"
	MaterialPodocarpus      Material = "Podocarpus"
	MaterialCedrela         Material = "Cedrela"
	MaterialEntandrophragma Material = "Entandrophragma"
	MaterialMahogany        Material = "Libuyu"
)

// User Model
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FirstName string    `gorm:"not null"`
	OtherName string    `gorm:"not null"`
	Email     string    `gorm:"unique;not null"`
	Phone     string    `gorm:"not null;default:'0781234567'"`
	Password  string    `gorm:"not null"`
	Salt      string    `gorm:"not null"`
	Role      Role      `gorm:"type:varchar(20);default:'CUSTOMER'"`

	// Location Fields
	City    string `gorm:"not null;default:'Kigali'"`
	Address string `gorm:"not null;default:'KN 123 ST'"`

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
	Owner    Vendor    `gorm:"foreignKey:VendorID;constraint:OnDelete:CASCADE"`
	Name     string    `gorm:"not null"`
	Price    float64   `gorm:"not null"`
	Category string    `gorm:"not null;type:varchar(30)"`
	Material Material  `gorm:"not null"`
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

	PaymentStatus PaymentStatus `gorm:"type:varchar(20);default:'PENDING'"`
	PaymentRef    string        `gorm:"type:varchar(100);"`
}

// Wishlist Model
type Wishlist struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	ProductID uuid.UUID `gorm:"not null"`
	Product   Product   `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"`
}

type Payment struct {
	ID            uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID       uuid.UUID     `gorm:"not null"`
	Order         Order         `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	Amount        float64       `gorm:"not null"`
	Method        PaymentMethod `gorm:"type:varchar(20);not null"`
	Status        PaymentStatus `gorm:"type:varchar(20);default:'PENDING'"`
	TransactionID string        `gorm:"type:varchar(100);unique"`
	CreatedAt     time.Time
}
