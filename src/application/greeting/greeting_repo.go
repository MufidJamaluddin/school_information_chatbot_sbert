package greeting

import (
	"chatbot_be_go/src/application/greeting/dto"
	"context"
)

type IGreetingRepository interface {
	FindCurrentGreeting(ctx context.Context) (greeting string, err error)

	ListGreeting(
		ctx context.Context,
		keyword string,
		start uint,
		size uint,
		callbackSendData func(adminItem interface{}),
	) (totalAll uint64, err error)

	SaveNewGreeting(
		ctx context.Context,
		greeting *dto.CreateGreetingDTO,
	) (uint64, error)

	UpdateGreeting(
		ctx context.Context,
		greetingId uint64,
		greeting *dto.UpdateGreetingDTO,
	) error

	DeleteGreeting(
		ctx context.Context,
		greetingId uint64,
	) error
}
