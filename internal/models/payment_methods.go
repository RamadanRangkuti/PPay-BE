package models

import (
	"time"
)

type PaymentMethod struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Tax       float64   `gorm:"type:money;default:0.00"`
	IsDeleted bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	TopupTransactions []TopupTransaction `gorm:"constraint:OnDelete:CASCADE"`
}
