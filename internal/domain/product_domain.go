package domain

type Product struct {
	Base
	Name        string      `json:"name" gorm:"size:255;not null"`
	Description string      `json:"description" gorm:"type:text"`
	Price       float64     `json:"price" gorm:"type:decimal(10,2);not null"`
	SKU         string      `json:"sku" gorm:"uniqueIndex;size:50;not null"`
	Stock       int         `json:"stock" gorm:"not null"`
	CategoryID  uint        `json:"category_id" gorm:"index"`
	Category    Category    `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	ImageURL    string      `json:"image_url" gorm:"size:255"`
	OrderItems  []OrderItem `json:"-" gorm:"foreignKey:ProductID"`
}

type Category struct {
	Base
	Name        string    `json:"name" gorm:"size:100;not null"`
	Description string    `json:"description" gorm:"type:text"`
	Products    []Product `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}
