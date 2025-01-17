package dto

type UserSummaryDTO struct {
	Image    *string `json:"image"`
	Fullname string  `json:"fullname"`
	Phone    *string `json:"phone"`
}

type UpdateUserRequest struct {
	Fullname *string `json:"fullname" formL:"fullname"`
	Email    *string `json:"email" form:"email" binding:"omitempty,email"`
	Password *string `json:"password" form:"password" binding:"omitempty,min=6"`
	Pin      *string `json:"pin" form:"pin"`
	Phone    *string `json:"phone" form:"phone"`
	Image    *string `json:"image" form:"image"`
}
