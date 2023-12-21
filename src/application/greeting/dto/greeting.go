package dto

import "time"

type CreateGreetingDTO struct {
	Greeting  string `json:"greeting"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	CreatedBy string `json:"-"`
}

type UpdateGreetingDTO struct {
	Greeting  string `json:"greeting"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	UpdateBy  string `json:"-"`
}

type GreetingItemDTO struct {
	GreetingId uint64     `json:"id"`
	Greeting   string     `json:"greeting"`
	StartTime  string     `json:"startTime"`
	EndTime    string     `json:"endTime"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
	CreatedBy  string     `json:"createdBy"`
	UpdatedBy  *string    `json:"updatedBy"`
}

type ResponseMessageDTO struct {
	Id      uint64 `json:"id"`
	Message string `json:"message"`
}
