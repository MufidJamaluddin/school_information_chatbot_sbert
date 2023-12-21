package dto

import "time"

type MessageDTO struct {
	FromId       string
	FromPhoneNo  string
	FromFullName string
	TextMessage  string
	ChatAt       time.Time
}
