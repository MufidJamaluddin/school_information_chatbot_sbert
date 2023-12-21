package dto

import "chatbot_be_go/src/application/shared"

type LoginDTO struct {
	Username string `json:"username" validate:"nonzero,max=10"`
	Password string `json:"password" validate:"nonzero,max=60"`
}

type AuthResponse struct {
	Token    string           `json:"token"`
	UserData *shared.UserData `json:"userData"`
}
