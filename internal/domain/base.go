package domain

import (
	"time"
)

type Base struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
