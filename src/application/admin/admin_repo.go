package admin

import (
	"chatbot_be_go/src/application/admin/dto"
	"context"
)

type IAdminRepository interface {
	CreateNewAdmin(
		ctx context.Context,
		admin dto.Admin,
		createdByUserName string,
	) error

	UpdateAdmin(
		ctx context.Context,
		userName string,
		admin dto.Admin,
		updatedByUserName string,
	) error

	ListAdmin(
		ctx context.Context,
		keyword string,
		start uint,
		size uint,
		callbackSendData func(adminItem interface{}),
	) (totalAll uint64, err error)
}
