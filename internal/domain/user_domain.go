package domain

type User struct {
	Base
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"`
	FirstName string    `json:"first_name" gorm:"size:100"`
	LastName  string    `json:"last_name" gorm:"size:100"`
	Role      UserRole  `json:"role" gorm:"type:varchar(20);default:'user'"`
	Orders    []Order   `json:"orders,omitempty" gorm:"foreignKey:UserID"`
	Addresses []Address `json:"addresses,omitempty" gorm:"foreignKey:UserID"`
}

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)
