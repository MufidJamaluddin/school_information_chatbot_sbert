package postgres

import (
	"chatbot_be_go/src/application/admin"
	"chatbot_be_go/src/application/admin/dto"
	"context"
	"fmt"
)

type adminRepo struct {
	db IDB
}

var _ admin.IAdminRepository = &adminRepo{}

func NewAdminRepository(db IDB) admin.IAdminRepository {
	return &adminRepo{
		db: db,
	}
}

func (a *adminRepo) CreateNewAdmin(
	ctx context.Context,
	admin dto.Admin,
	createdByUserName string,
) error {
	sqlDb := a.db.GetSqlDb()

	_, err := sqlDb.ExecContext(
		ctx,
		"INSERT INTO public.\"admin\"(username, password, full_name, phone_no, position, role_group_id, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7);",
		&admin.Username,
		&admin.HashedPassword,
		&admin.FullName,
		&admin.PhoneNo,
		&admin.Position,
		&admin.RoleGroupId,
		&createdByUserName,
	)

	return err
}

func (a *adminRepo) UpdateAdmin(
	ctx context.Context,
	username string,
	admin dto.Admin,
	updatedBy string,
) error {
	sqlDb := a.db.GetSqlDb()

	_, err := sqlDb.ExecContext(
		ctx,
		"UPDATE public.\"admin\" SET password = $1, full_name = $2, phone_no = $3, position = $4, role_group_id = $5, updated_by = $6, updated_at = NOW() WHERE username = $7;",
		&admin.HashedPassword,
		&admin.FullName,
		&admin.PhoneNo,
		&admin.Position,
		&admin.RoleGroupId,
		&username,
		&updatedBy,
	)

	return err
}

func (a *adminRepo) ListAdmin(
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

	var adminItem dto.ListAdminItem

	sqlDb := a.db.GetSqlDb()

	if keyword != "" {
		iLikeKeyword := fmt.Sprintf("%%%s%", keyword)
		totalQuery = "SELECT count(*) FROM \"admin\" WHERE full_name ILIKE $1"
		listQuery = `SELECT
				username,
				full_name,
				phone_no,
				position,
				role_group_id
			FROM
				"admin"
			WHERE
				full_name ILIKE $1
			LIMIT $2
			OFFSET $3;`

		totalQueryArgs = append(totalQueryArgs, iLikeKeyword)
		listQueryArgs = append(listQueryArgs, iLikeKeyword)
	} else {
		totalQuery = "SELECT count(*) FROM \"admin\""
		listQuery = `SELECT
				username,
				full_name,
				phone_no,
				position,
				role_group_id
			FROM
				"admin"
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
			&adminItem.Username,
			&adminItem.FullName,
			&adminItem.PhoneNo,
			&adminItem.Position,
			&adminItem.RoleGroupId,
		)

		callbackSendData(&adminItem)
	}

	return
}
