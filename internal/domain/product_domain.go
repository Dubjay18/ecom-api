package domain

type Product struct {
	Base
	Name        string      `json:"name" gorm:"size:255;not null"`
	Description string      `json:"description" gorm:"type:text"`
	Price       float64     `json:"price" gorm:"type:decimal(10,2);not null"`
	SKU         string      `json:"sku" gorm:"uniqueIndex;size:50;not null"`
	Stock       int         `json:"stock" gorm:"not null"`
	Category    string      `json:"category" gorm:"type:varchar(100);"`
	ImageURL    string      `json:"image_url" gorm:"size:255"`
	OrderItems  []OrderItem `json:"-" gorm:"foreignKey:ProductID"`
}

type ProductFilter struct {
	Name string

	MinPrice float64

	MaxPrice float64
}

type CreateProductRequest struct {
	Name        string  `form:"name" binding:"required"`
	Price       float64 `form:"price" binding:"required,gt=0"`
	Description string  `json:"description"`
	Stock       int     `form:"stock" binding:"required,gt=0"`
	SKU         string  `form:"sku" binding:"required"`
	Category    string  `form:"category"`
}

type UpdateProductRequest struct {
	Name        string  `form:"name"`
	Price       float64 `form:"price" binding:"gt=0"`
	Description string  `json:"description"`
	Stock       int     `form:"stock" binding:"gt=0"`
	SKU         string  `form:"sku"`
	Category    string  `form:"category"`
}
