package postgres

import (
	"chatbot_be_go/src/application/role_group"
	"chatbot_be_go/src/application/role_group/dto"
	"context"
	"fmt"
)

type roleGroupRepository struct {
	db IDB
}

var _ role_group.IRoleGroupRepository = &roleGroupRepository{}

func NewRoleGroupRepository(db IDB) role_group.IRoleGroupRepository {
	return &roleGroupRepository{
		db: db,
	}
}

func (r *roleGroupRepository) ListRoleGroup(
	ctx context.Context,
	keyword string,
	start uint,
	size uint,
	callbackSendData func(roleItem interface{}),
) (totalAll uint64, err error) {
	var totalQuery string
	var listQuery string

	var totalQueryArgs []interface{}
	var listQueryArgs []interface{}

	var roleGroupItemDto dto.RoleItemGroupDTO

	sqlDb := r.db.GetSqlDb()

	if keyword != "" {
		iLikeKeyword := fmt.Sprintf("%%%s%%", keyword)
		totalQuery = "SELECT COUNT(*) FROM public.\"role_group\" WHERE role_group ILIKE $1;"
		listQuery = `
			SELECT
				id,
				role_group,
				created_at,
				updated_at,
				created_by,
				updated_by
			FROM
				public."role_group"
			WHERE
				role_group ILIKE $1
			LIMIT $2 
			OFFSET $3;`

		totalQueryArgs = append(totalQueryArgs, iLikeKeyword)
		listQueryArgs = append(listQueryArgs, iLikeKeyword)
	} else {
		totalQuery = "SELECT COUNT(*) FROM public.\"role_group\";"
		listQuery = `
			SELECT
				id,
				role_group,
				created_at,
				updated_at,
				created_by,
				updated_by
			FROM
				public."role_group"
			LIMIT $1 
			OFFSET $2;`
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
			&roleGroupItemDto.RoleGroupId,
			&roleGroupItemDto.RoleGroup,
			&roleGroupItemDto.CreatedAt,
			&roleGroupItemDto.UpdatedAt,
			&roleGroupItemDto.CreatedBy,
			&roleGroupItemDto.UpdatedBy,
		)

		callbackSendData(&roleGroupItemDto)
	}

	return
}

func (r *roleGroupRepository) SaveNewRoleGroup(
	ctx context.Context,
	roleGroup *dto.CreateRoleGroupDTO,
) (id uint64, err error) {
	sqlDb := r.db.GetSqlDb()

	err = sqlDb.QueryRowContext(
		ctx,
		"INSERT INTO public.\"role_group\" (role_group, created_by, created_at) VALUES ($1, $2, NOW()) RETURNING id;",
		&roleGroup.RoleGroup,
		&roleGroup.CreatedBy,
	).Scan(&id)

	return
}

func (r *roleGroupRepository) UpdateRoleGroup(
	ctx context.Context,
	roleGroupId uint64,
	roleGroup *dto.UpdateRoleGroupDTO,
) error {
	sqlDb := r.db.GetSqlDb()

	_, err := sqlDb.ExecContext(
		ctx,
		"UPDATE public.\"role_group\" SET role_group = $1, updated_by = $2, updated_at = NOW() WHERE id = $3;",
		&roleGroup.RoleGroup,
		&roleGroup.UpdatedBy,
		&roleGroupId,
	)

	return err
}
