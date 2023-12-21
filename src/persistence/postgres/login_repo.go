package postgres

import (
	"chatbot_be_go/src/application/login"
	"chatbot_be_go/src/application/shared"
	"context"
)

type loginRepository struct {
	db IDB
}

var _ login.ILoginRepository = &loginRepository{}

func NewLoginRepository(db IDB) login.ILoginRepository {
	return &loginRepository{
		db: db,
	}
}

func (l *loginRepository) GetUserData(ctx context.Context, userName string) (resp shared.UserData, err error) {
	sqlDb := l.db.GetSqlDb()
	resp.UserName = userName

	err = sqlDb.QueryRowContext(
		ctx,
		`SELECT
			role_group_id,
			full_name,
			position,
			password
		FROM
			public.admin
		WHERE
			username = $1;`,
		&userName,
	).Scan(
		&resp.RoleGroupId,
		&resp.FullName,
		&resp.Position,
		&resp.PasswordHash,
	)

	return
}
