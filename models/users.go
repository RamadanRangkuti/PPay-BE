package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Fullname  string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"type:varchar(255);unique;not null"`
	Password  string    `gorm:"type:varchar(255);not null"`
	Pin       *string   `gorm:"type:char(6)"`
	Phone     *string   `gorm:"type:varchar(16);unique"`
	Image     *string   `gorm:"type:varchar(255)"`
	IsDeleted bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Wallet       Wallet        `gorm:"constraint:OnDelete:CASCADE"`
	Transactions []Transaction `gorm:"constraint:OnDelete:CASCADE"`
}
