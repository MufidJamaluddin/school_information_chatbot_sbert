package chat_user

import "context"

type IChatUserRepository interface {
	ListChatUser(
		ctx context.Context,
		keyword string,
		start uint,
		size uint,
		callbackSendData func(adminItem interface{}),
	) (totalAll uint64, err error)
}
