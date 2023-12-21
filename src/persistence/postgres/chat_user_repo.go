package postgres

import (
	"chatbot_be_go/src/application/chat_user"
	"chatbot_be_go/src/application/chat_user/dto"
	"context"
	"fmt"
)

type chatUserRepository struct {
	db IDB
}

var _ chat_user.IChatUserRepository = &chatUserRepository{}

func NewChatUserRepository(db IDB) chat_user.IChatUserRepository {
	return &chatUserRepository{
		db: db,
	}
}

func (c chatUserRepository) ListChatUser(
	ctx context.Context,
	keyword string,
	start uint,
	size uint,
	callbackSendData func(adminItem interface{}),
) (totalAll uint64, err error) {
	var totalQuery string
	var listQuery string

	var totalQueryArgs []interface{}
	var listQueryArgs []interface{}

	var chatUserItem dto.ChatUserItemDTO

	sqlDb := c.db.GetSqlDb()

	if keyword != "" {
		iLikeKeyword := fmt.Sprintf("%%%s%", keyword)
		totalQuery = "SELECT count(*) FROM public.\"chat_user\" ILIKE $1"
		listQuery = "SELECT phone_no, full_name FROM public.\"chat_user\" WHERE full_name ILIKE $1 LIMIT $2 OFFSET $3;"

		totalQueryArgs = append(totalQueryArgs, iLikeKeyword)
		listQueryArgs = append(listQueryArgs, iLikeKeyword)
	} else {
		totalQuery = "SELECT count(*) FROM public.\"chat_user\";"
		listQuery = "SELECT phone_no, full_name FROM public.\"chat_user\" LIMIT $1 OFFSET $2;"
	}

	listQueryArgs = append(listQueryArgs, size)
	listQueryArgs = append(listQueryArgs, start)

	if err = sqlDb.QueryRowContext(
		ctx,
		totalQuery,
		totalQueryArgs...,
	).Scan(&totalAll); err != nil {
		return
	}

	rows, err := sqlDb.QueryContext(
		ctx,
		listQuery,
		listQueryArgs...,
	)
	if err != nil {
		return
	}

	for rows.Next() {
		_ = rows.Scan(
			&chatUserItem.PhoneNo,
			&chatUserItem.FullName,
		)

		callbackSendData(&chatUserItem)
	}

	return
}
