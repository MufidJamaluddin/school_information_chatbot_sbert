package postgres

import (
	uq "chatbot_be_go/src/application/unanswered_question"
	"chatbot_be_go/src/application/unanswered_question/dto"
	"context"
	"fmt"
)

type unansweredQuestionRepository struct {
	db IDB
}

var _ uq.IUnansweredQuestionRepository = &unansweredQuestionRepository{}

func NewUnansweredQuestionRepository(
	db IDB,
) uq.IUnansweredQuestionRepository {
	return &unansweredQuestionRepository{
		db: db,
	}
}

func (u *unansweredQuestionRepository) ListUnansweredQuestion(
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

	var unansweredItemDto dto.UnansweredQuestionDTO

	sqlDb := u.db.GetSqlDb()

	if keyword != "" {
		iLikeKeyword := fmt.Sprintf("%%%s%", keyword)
		totalQuery = `
			SELECT
				COUNT(*)
			FROM
				public."user_question" uq
			WHERE
				uq.should_answer_manually IS TRUE
				AND
				uq.manual_answer_by IS NULL
				AND
			  uq.question ILIKE $1;`

		listQuery = `
			SELECT
				uq.id,
				uq.phone_no,
				cu.full_name,
				uq.question
			FROM
				public."user_question" uq
			INNER JOIN
				public."chat_user" cu
				ON cu.phone_no = uq.phone_no
			WHERE
				uq.should_answer_manually IS TRUE
				AND
				uq.manual_answer_by IS NULL
				AND
			  uq.question ILIKE $1
			LIMIT $2 
			OFFSET $3;`

		totalQueryArgs = append(totalQueryArgs, iLikeKeyword)
		listQueryArgs = append(listQueryArgs, iLikeKeyword)
	} else {
		totalQuery = `
			SELECT
				COUNT(*)
			FROM
				public."user_question" uq
			WHERE
				uq.should_answer_manually IS TRUE
				AND
				uq.manual_answer_by IS NULL;`

		listQuery = `
			SELECT
				uq.id,
				uq.phone_no,
				cu.full_name,
				uq.question
			FROM
				public."user_question" uq
			INNER JOIN
				public."chat_user" cu
				ON cu.phone_no = uq.phone_no
			WHERE
				uq.should_answer_manually IS TRUE
				AND
				uq.manual_answer_by IS NULL
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
			&unansweredItemDto.Id,
			&unansweredItemDto.PhoneNo,
			&unansweredItemDto.FullName,
			&unansweredItemDto.Question,
		)

		callbackSendData(&unansweredItemDto)
	}

	return
}

func (u *unansweredQuestionRepository) AnswerQuestion(
	ctx context.Context,
	userQuestionId uint64,
	answer string,
	answerByUserName string,
) (err error) {
	sqlDb := u.db.GetSqlDb()

	_, err = sqlDb.ExecContext(
		ctx,
		"UPDATE public.\"user_question\" SET answer = $1, manual_answer_by = $2, answer_at = NOW() WHERE id = $3;",
		answer,
		answerByUserName,
		userQuestionId,
	)

	return
}
