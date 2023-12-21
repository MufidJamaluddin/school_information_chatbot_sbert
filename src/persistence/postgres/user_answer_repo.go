package postgres

import (
	ua "chatbot_be_go/src/application/user_response"
	"chatbot_be_go/src/application/user_response/dto"
	"context"
)

type userResponseRepository struct {
	db IDB
}

func NewUserResponseRepository(db IDB) ua.IUserResponseRepository {
	return &userResponseRepository{
		db: db,
	}
}

var _ ua.IUserResponseRepository = &userResponseRepository{}

func (r *userResponseRepository) CreateNewUserResponse(
	ctx context.Context,
	userAnswer *dto.UserResponse,
) (err error) {
	sqlDb := r.db.GetSqlDb()

	_, err = sqlDb.ExecContext(
		ctx,
		"INSERT INTO public.\"user_response\" (user_id, question, answer, score) VALUES ($1, $2, $3, $4);",
		userAnswer.UserId,
		userAnswer.Question,
		userAnswer.Answer,
		userAnswer.Score,
	)

	return
}
