package postgres

import (
	fe "chatbot_be_go/src/application/question"
	"chatbot_be_go/src/application/question/dto"
	dm "chatbot_be_go/src/domain"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type questionRepository struct {
	db         IDB
	vectorizer dm.ISBertVectorizer
}

var _ fe.IQuestionRepository = &questionRepository{}

func NewQuestionRepository(
	db IDB,
	vectorizer dm.ISBertVectorizer,
) fe.IQuestionRepository {
	return &questionRepository{
		db:         db,
		vectorizer: vectorizer,
	}
}

func (t *questionRepository) FindAnswer(nearestAnswer string) (answer string, err error) {
	sqlDb := t.db.GetSqlDb()
	var similarityValue float64

	sbertVector, err := t.vectorizer.Encode(nearestAnswer)
	if err != nil {
		return
	}

	params := strings.ReplaceAll(
		fmt.Sprintf("%f", sbertVector),
		" ",
		",",
	)

	err = sqlDb.QueryRow(
		fmt.Sprintf(
			`SELECT
				answer,
				question_vector_sbert <-> '%s' AS similarity
			FROM
				"question"
			WHERE
				similarity > 0
			ORDER BY
				similarity
			LIMIT 1;`,
			params,
		),
	).Scan(&answer, &similarityValue)

	if similarityValue < 0 {
		err = sql.ErrNoRows
		return
	}

	return
}

func (t *questionRepository) FindAnswerWithSimilarityValue(
	ctx context.Context,
	nearestAnswer string,
) (
	answer string,
	similarityValue float64,
	err error,
) {
	sqlDb := t.db.GetSqlDb()

	sbertVector, err := t.vectorizer.Encode(nearestAnswer)
	if err != nil {
		return
	}

	params := strings.ReplaceAll(
		fmt.Sprintf("%f", sbertVector),
		" ",
		",",
	)

	err = sqlDb.QueryRowContext(
		ctx,
		fmt.Sprintf(
			`SELECT
				answer,
				question_vector_sbert <-> '%s' AS similarity
			FROM
				"question"
			WHERE
				similarity > 0
			ORDER BY
				similarity
			LIMIT 1;`,
			params,
		),
	).Scan(&answer, &similarityValue)

	return
}

func (t *questionRepository) ListQuestion(
	ctx context.Context,
	keyword string,
	start uint,
	size uint,
	callbackSendData func(questionItem interface{}),
) (totalAll uint64, err error) {
	var totalQuery string
	var listQuery string

	var totalQueryArgs []interface{}
	var listQueryArgs []interface{}

	var questionItemDto dto.QuestionItemDTO

	sqlDb := t.db.GetSqlDb()

	if keyword != "" {
		iLikeKeyword := fmt.Sprintf("%%%s%%", keyword)
		totalQuery = "SELECT COUNT(*) FROM public.\"question\" WHERE question ILIKE $1"
		listQuery = `
			SELECT
				id,
				question,
				answer,
				role_group_id,
				created_by,
				updated_by
			FROM
				public."question"
			WHERE 
			  question ILIKE $1
			LIMIT $2 
			OFFSET $3;`

		totalQueryArgs = append(totalQueryArgs, iLikeKeyword)
		listQueryArgs = append(listQueryArgs, iLikeKeyword)
	} else {
		totalQuery = "SELECT COUNT(*) FROM public.\"question\";"
		listQuery = `
			SELECT
				id,
				question,
				answer,
				role_group_id,
				created_by,
				updated_by
			FROM
				public."question"
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
			&questionItemDto.QuestionId,
			&questionItemDto.Question,
			&questionItemDto.Answer,
			&questionItemDto.RoleGroupId,
			&questionItemDto.CreatedBy,
			&questionItemDto.UpdatedBy,
		)

		callbackSendData(&questionItemDto)
	}

	return
}

func (t *questionRepository) SaveNewQuestion(
	ctx context.Context,
	question string,
	answer string,
	roleGroupId uint64,
	createdBy string,
) (questionId uint64, err error) {
	sbertVector, err := t.vectorizer.Encode(question)
	if err != nil {
		return
	}

	sbertVectorStr := strings.ReplaceAll(
		fmt.Sprintf("%f", sbertVector),
		" ",
		",",
	)

	trx, err := t.db.GetSqlDb().BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

	if err != nil {
		return
	}

	stmt, err := trx.Prepare(
		`INSERT INTO "question" (role_group_id, question, question_vector_sbert, answer, created_by) VALUES ($1, $2, $3, $4, $5) RETURNING id;`,
	)

	if err != nil {
		_ = trx.Rollback()
		return
	}

	if err = stmt.QueryRow(
		roleGroupId,
		question,
		sbertVectorStr,
		answer,
		createdBy,
	).Scan(&questionId); err != nil {
		_ = trx.Rollback()
		return
	}

	err = trx.Commit()
	return
}

func (t *questionRepository) UpdateQuestion(
	ctx context.Context,
	questionId uint64,
	question string,
	answer string,
	updatedBy string,
) (err error) {
	sbertVector, err := t.vectorizer.Encode(question)
	if err != nil {
		return
	}

	sbertVectorStr := strings.ReplaceAll(
		fmt.Sprintf("%f", sbertVector),
		" ",
		",",
	)

	trx, err := t.db.GetSqlDb().BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

	if err != nil {
		return
	}

	if _, err = trx.Exec(
		`UPDATE "question" SET question = $1, question_vector_sbert=$2, answer = $3, updated_by = $4 WHERE id = $5;`,
		question,
		sbertVectorStr,
		answer,
		updatedBy,
		questionId,
	); err != nil {
		_ = trx.Rollback()
		return
	}

	err = trx.Commit()
	return
}

func (t *questionRepository) DeleteQuestion(
	ctx context.Context,
	questionId uint64,
) (err error) {
	trx, err := t.db.GetSqlDb().BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

	if err != nil {
		return
	}

	if _, err = trx.Exec(
		`DELETE FROM public."question" WHERE id = $1;`,
		questionId,
	); err != nil {
		_ = trx.Rollback()
		return
	}

	err = trx.Commit()
	return
}

func (t *questionRepository) TruncateQuestion(ctx context.Context) (err error) {
	trx, err := t.db.GetSqlDb().BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelSerializable,
	})

	if err != nil {
		return
	}

	if _, err = trx.Exec(`TRUNCATE TABLE question RESTART IDENTITY CASCADE;`); err != nil {
		_ = trx.Rollback()
		return
	}

	err = trx.Commit()
	return
}
