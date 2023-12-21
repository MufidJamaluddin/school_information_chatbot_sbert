package login

import (
	"chatbot_be_go/src/application/shared"
	"context"
)

type ILoginRepository interface {
	GetUserData(ctx context.Context, userName string) (resp shared.UserData, err error)
}
