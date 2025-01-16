package models

import (
	"time"
)

type TopupTransaction struct {
	ID              uint      `gorm:"primaryKey"`
	TransactionID   uint      `gorm:"not null;constraint:OnDelete:CASCADE"`
	PaymentMethodID uint      `gorm:"not null;constraint:OnDelete:CASCADE"`
	IsDeleted       bool      `gorm:"default:false"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}
