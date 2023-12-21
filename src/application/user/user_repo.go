package admin

import (
	"chatbot_be_go/src/application/user/dto"
	"context"
)

type IUserRepository interface {
	CreateNewUser(
		ctx context.Context,
		admin *dto.User,
	) error
}
