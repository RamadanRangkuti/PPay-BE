package dto

type UserSummaryDTO struct {
	Image    *string `json:"image"`
	Fullname string  `json:"fullname"`
	Phone    *string `json:"phone"`
}
