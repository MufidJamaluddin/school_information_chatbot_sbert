package postgres

import (
	ua "chatbot_be_go/src/application/user"
	"chatbot_be_go/src/application/user/dto"
	"context"
)

type userRepository struct {
	db IDB
}

func NewUserRepository(db IDB) ua.IUserRepository {
	return &userRepository{
		db: db,
	}
}

var _ ua.IUserRepository = &userRepository{}

func (r *userRepository) CreateNewUser(ctx context.Context, user *dto.User) (err error) {
	sqlDb := r.db.GetSqlDb()

	err = sqlDb.QueryRowContext(
		ctx,
		"INSERT INTO public.\"user\" (full_name, is_student, user_role, class_name, age) VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		user.FullName,
		user.GetIsStudent(),
		user.UserRole,
		user.ClassName,
		user.Age,
	).Scan(
		&user.Id,
	)

	return
}
