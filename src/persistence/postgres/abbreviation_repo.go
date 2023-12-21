package postgres

import (
	"chatbot_be_go/src/application/abbreviation"
	"chatbot_be_go/src/application/abbreviation/dto"
	"context"
	"fmt"

	"github.com/lib/pq"
)

type abbreviationRepository struct {
	db IDB
}

var _ abbreviation.IAbbreviationRepository = &abbreviationRepository{}

func NewAbbreviationRepository(db IDB) abbreviation.IAbbreviationRepository {
	return &abbreviationRepository{
		db: db,
	}
}

func (s *abbreviationRepository) GetAbbreviation(
	ctx context.Context,
	standardWord string,
) (*dto.AbbreviationItemDTO, error) {
	var err error
	var abbreviationItemDto dto.AbbreviationItemDTO

	sqlDb := s.db.GetSqlDb()

	err = sqlDb.QueryRowContext(
		ctx,
		`SELECT
			list_abbreviation_term,
			created_at,
			created_by,
			updated_at,
			updated_by
		FROM
			public."abbreviation"
		WHERE 
			standard_word = $1`,
		standardWord,
	).Scan(
		pq.Array(&abbreviationItemDto.ListAbbreviationTerm),
		&abbreviationItemDto.CreatedAt,
		&abbreviationItemDto.CreatedBy,
		&abbreviationItemDto.UpdatedAt,
		&abbreviationItemDto.UpdatedBy,
	)

	abbreviationItemDto.StandardWord = standardWord

	return &abbreviationItemDto, err
}

func (s *abbreviationRepository) ListAbbreviation(
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

	var abbreviationItemDto dto.AbbreviationItemDTO

	sqlDb := s.db.GetSqlDb()

	if keyword != "" {
		iLikeKeyword := fmt.Sprintf("%%%s%%", keyword)
		totalQuery = "SELECT COUNT(*) FROM public.\"abbreviation\" WHERE standard_word ILIKE $1 OR $2 >= ANY (list_abbreviation_term);"
		listQuery = `
			SELECT
				standard_word,
				list_abbreviation_term,
				created_at,
				created_by,
				updated_at,
				updated_by
			FROM
				public."abbreviation"
			WHERE 
			  standard_word ILIKE $1 OR $2 >= ANY (list_abbreviation_term)
			LIMIT $3 
			OFFSET $4;`

		totalQueryArgs = append(totalQueryArgs, iLikeKeyword)
		totalQueryArgs = append(totalQueryArgs, keyword)

		listQueryArgs = append(listQueryArgs, iLikeKeyword)
		listQueryArgs = append(listQueryArgs, keyword)
	} else {
		totalQuery = "SELECT COUNT(*) FROM public.\"abbreviation\";"
		listQuery = `
			SELECT
				standard_word,
				list_abbreviation_term,
				created_at,
				created_by,
				updated_at,
				updated_by
			FROM
				public."abbreviation"
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
			&abbreviationItemDto.StandardWord,
			pq.Array(&abbreviationItemDto.ListAbbreviationTerm),
			&abbreviationItemDto.CreatedAt,
			&abbreviationItemDto.CreatedBy,
			&abbreviationItemDto.UpdatedAt,
			&abbreviationItemDto.UpdatedBy,
		)

		callbackSendData(&abbreviationItemDto)
	}

	return
}

func (s *abbreviationRepository) SaveNewAbbreviation(
	ctx context.Context,
	abbreviation *dto.CreateAbbreviationDTO,
) (err error) {
	sqlDb := s.db.GetSqlDb()

	_, err = sqlDb.ExecContext(
		ctx,
		"INSERT INTO public.\"abbreviation\" (standard_word, list_abbreviation_term, created_by, created_at) VALUES ($1, $2, $3, NOW());",
		&abbreviation.StandardWord,
		pq.Array(&abbreviation.ListAbbreviationTerm),
		&abbreviation.CreatedBy,
	)

	return
}

func (s *abbreviationRepository) UpdateAbbreviation(
	ctx context.Context,
	abbreviation *dto.UpdateAbbreviationDTO,
) error {
	sqlDb := s.db.GetSqlDb()

	_, err := sqlDb.ExecContext(
		ctx,
		"UPDATE public.\"abbreviation\" SET list_abbreviation_term = $1, updated_by = $2, updated_at = NOW() WHERE standard_word = $3;",
		pq.Array(&abbreviation.ListAbbreviationTerm),
		&abbreviation.UpdatedBy,
		&abbreviation.StandardWord,
	)

	return err
}

func (s *abbreviationRepository) DeleteAbbreviation(
	ctx context.Context,
	standardWord string,
) error {
	sqlDb := s.db.GetSqlDb()

	_, err := sqlDb.ExecContext(
		ctx,
		"DELETE FROM public.\"abbreviation\" WHERE standard_word = $1;",
		&standardWord,
	)

	return err
}
