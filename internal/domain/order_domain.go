package domain

type Order struct {
	Base
	UserID            uint        `json:"user_id" gorm:"index;not null"`
	User              User        `json:"-" gorm:"foreignKey:UserID"`
	Status            OrderStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
	TotalAmount       float64     `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	Items             []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	ShippingAddressID uint        `json:"shipping_address_id" gorm:"not null"`
	// ShippingAddr   Address       `json:"address,omitempty" gorm:"foreignKey:ShippingAddrID"`
	PaymentStatus PaymentStatus `json:"payment_status" gorm:"type:varchar(20);default:'pending'"`
}

type OrderItem struct {
	Base
	OrderID   uint    `json:"-" gorm:"index;not null"`
	ProductID uint    `json:"-" gorm:"index;not null"`
	Product   Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Price     float64 `json:"price" gorm:"type:decimal(10,2);not null"`
}

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusConfirmed OrderStatus = "confirmed"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
	StatusCancelled OrderStatus = "cancelled"
)

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentCompleted PaymentStatus = "completed"
	PaymentFailed    PaymentStatus = "failed"
	PaymentRefunded  PaymentStatus = "refunded"
)

type Address struct {
	Base
	UserID     uint   `json:"user_id" gorm:"index"`
	Street     string `json:"street" binding:"required" gorm:"size:255;not null"`
	City       string `json:"city" binding:"required" gorm:"size:100;not null"`
	State      string `json:"state" binding:"required" gorm:"size:100;not null"`
	Country    string `json:"country" binding:"required" gorm:"size:100;not null"`
	PostalCode string `json:"postal_code" binding:"required" gorm:"size:20;not null"`
	IsDefault  bool   `json:"is_default" gorm:"default:false"`
}

type CreateOrderRequest struct {
	Items         []CreateOrderItem   `json:"items" binding:"required"`
	ShippingAddr  CreatAddressRequest `json:"shipping_address" binding:"required"`
	PaymentMethod string              `json:"payment_method" binding:"required"`
}

type CreateOrderItem struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status"`
}

type CreatAddressRequest struct {
	Street     string `json:"street" binding:"required"`
	City       string `json:"city" binding:"required"`
	State      string `json:"state" binding:"required"`
	Country    string `json:"country" binding:"required"`
	PostalCode string `json:"postal_code" binding:"required"`
	IsDefault  bool   `json:"is_default"`
}
