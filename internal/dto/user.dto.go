package dto

import "time"

type UserSummaryDTO struct {
	Id       uint    `json:"id"`
	Image    *string `json:"image"`
	Email    string  `json:"email,omitempty"`
	Fullname string  `json:"fullname"`
	Phone    *string `json:"phone"`
}

type UpdateUserRequest struct {
	Fullname  *string   `json:"fullname" form:"fullname"`
	Email     *string   `json:"email" form:"email" binding:"omitempty,email"`
	Password  *string   `json:"-" form:"password" binding:"omitempty,min=6"`
	Pin       *string   `json:"pin" form:"pin" binding:"omitempty"`
	Phone     *string   `json:"phone" form:"phone"`
	Image     *string   `json:"image" form:"image"`
	UpdatedAt time.Time `json:"updatedAt" form:"updatedAt"`
}

type CreateUserRequest struct {
	Fullname *string `json:"fullname"`
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Pin      *string `json:"pin"`
	Phone    *string `json:"phone"`
	Image    *string `json:"image"`
}
