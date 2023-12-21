package admin

import (
	"chatbot_be_go/src/application/user_response/dto"
	"context"
)

type IUserResponseRepository interface {
	CreateNewUserResponse(
		ctx context.Context,
		admin *dto.UserResponse,
	) error
}
