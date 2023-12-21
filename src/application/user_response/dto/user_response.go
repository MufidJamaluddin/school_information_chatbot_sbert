package dto

type UserResponse struct {
	UserId   uint   `json:"userId" validate:"nonzero"`
	Question string `json:"question" validate:"nonzero"`
	Answer   string `json:"answer"`
	Score    uint8  `json:"score" validate:"min=0,max=5"`
}
