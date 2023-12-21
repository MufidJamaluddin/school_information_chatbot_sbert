package role_group

import (
	"chatbot_be_go/src/application/role_group/dto"
	"context"
)

type IRoleGroupRepository interface {
	ListRoleGroup(
		ctx context.Context,
		keyword string,
		start uint,
		size uint,
		callbackSendData func(roleItem interface{}),
	) (totalAll uint64, err error)

	SaveNewRoleGroup(
		ctx context.Context,
		role *dto.CreateRoleGroupDTO,
	) (uint64, error)

	UpdateRoleGroup(
		ctx context.Context,
		roleGroupId uint64,
		role *dto.UpdateRoleGroupDTO,
	) error
}
