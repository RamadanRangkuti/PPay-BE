package models

import (
	"time"
)

type Transaction struct {
	ID              uint      `gorm:"primaryKey"`
	UserID          uint      `gorm:"not null;constraint:OnDelete:CASCADE"`
	Amount          float64   `gorm:"type:money;not null"`
	TransactionType string    `gorm:"type:transaction_type_enum;not null"`
	Notes           *string   `gorm:"type:varchar(255)"`
	IsDeleted       bool      `gorm:"default:false"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`

	TransferTransactions []TransferTransaction `gorm:"constraint:OnDelete:CASCADE"`
	TopupTransactions    []TopupTransaction    `gorm:"constraint:OnDelete:CASCADE"`
}
