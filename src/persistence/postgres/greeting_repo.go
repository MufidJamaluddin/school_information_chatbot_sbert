package postgres

import (
	"chatbot_be_go/src/application/greeting"
	"chatbot_be_go/src/application/greeting/dto"
	"context"
	"fmt"
)

type greetingRepository struct {
	db IDB
}

var _ greeting.IGreetingRepository = &greetingRepository{}

func NewGreetingRepository(db IDB) greeting.IGreetingRepository {
	return &greetingRepository{
		db: db,
	}
}

func (g *greetingRepository) FindCurrentGreeting(ctx context.Context) (greeting string, err error) {
	sqlDb := g.db.GetSqlDb()

	err = sqlDb.QueryRowContext(
		ctx,
		`SELECT
			greeting
		FROM
			public."greeting"
		WHERE
			(NOW() AT TIME ZONE 'Asia/Jakarta')::time BETWEEN start_time AND end_time
		LIMIT 1;`,
	).Scan(&greeting)

	return
}

func (g *greetingRepository) ListGreeting(
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

	var greetingItemDto dto.GreetingItemDTO

	sqlDb := g.db.GetSqlDb()

	if keyword != "" {
		iLikeKeyword := fmt.Sprintf("%%%s%%", keyword)
		totalQuery = "SELECT COUNT(*) FROM public.\"greeting\" WHERE greeting ILIKE $1"
		listQuery = `
			SELECT
				id,
				greeting,
				created_at,
				updated_at,
				created_by,
				updated_by,
				start_time,
				end_time
			FROM
				public."greeting"
			WHERE 
			    greeting ILIKE $1
			LIMIT $2 
			OFFSET $3;`

		totalQueryArgs = append(totalQueryArgs, iLikeKeyword)
		listQueryArgs = append(listQueryArgs, iLikeKeyword)
	} else {
		totalQuery = "SELECT COUNT(*) FROM public.\"greeting\";"
		listQuery = `
			SELECT
				id,
				greeting,
				created_at,
				updated_at,
				created_by,
				updated_by,
				start_time,
				end_time
			FROM
				public."greeting"
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
			&greetingItemDto.GreetingId,
			&greetingItemDto.Greeting,
			&greetingItemDto.CreatedAt,
			&greetingItemDto.UpdatedAt,
			&greetingItemDto.CreatedBy,
			&greetingItemDto.UpdatedBy,
			&greetingItemDto.StartTime,
			&greetingItemDto.EndTime,
		)

		callbackSendData(&greetingItemDto)
	}

	return
}

func (g *greetingRepository) SaveNewGreeting(
	ctx context.Context,
	greeting *dto.CreateGreetingDTO,
) (id uint64, err error) {
	sqlDb := g.db.GetSqlDb()

	err = sqlDb.QueryRowContext(
		ctx,
		"INSERT INTO public.\"greeting\" (greeting, start_time, end_time, created_by, created_at) VALUES ($1, $2, $3, $4, NOW()) RETURNING id;",
		&greeting.Greeting,
		&greeting.StartTime,
		&greeting.EndTime,
		&greeting.CreatedBy,
	).Scan(&id)

	return
}

func (g *greetingRepository) UpdateGreeting(
	ctx context.Context,
	greetingId uint64,
	greeting *dto.UpdateGreetingDTO,
) error {
	sqlDb := g.db.GetSqlDb()

	_, err := sqlDb.ExecContext(
		ctx,
		"UPDATE public.\"greeting\" SET greeting = $1, start_time = $2, end_time = $3, updated_by = $4, updated_at = NOW() WHERE id = $5;",
		&greeting.Greeting,
		&greeting.StartTime,
		&greeting.EndTime,
		&greeting.UpdateBy,
		&greetingId,
	)

	return err
}

func (g *greetingRepository) DeleteGreeting(
	ctx context.Context,
	greetingId uint64,
) error {
	sqlDb := g.db.GetSqlDb()

	_, err := sqlDb.ExecContext(
		ctx,
		"DELETE FROM public.\"greeting\" WHERE id = $1;",
		&greetingId,
	)

	return err
}
