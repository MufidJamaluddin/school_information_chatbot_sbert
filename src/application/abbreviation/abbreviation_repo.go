package abbreviation

import (
	"chatbot_be_go/src/application/abbreviation/dto"
	"context"
)

type IAbbreviationRepository interface {
	GetAbbreviation(
		ctx context.Context,
		standardWord string,
	) (*dto.AbbreviationItemDTO, error)

	ListAbbreviation(
		ctx context.Context,
		keyword string,
		start uint,
		size uint,
		callbackSendData func(roleItem interface{}),
	) (totalAll uint64, err error)

	SaveNewAbbreviation(
		ctx context.Context,
		abbreviation *dto.CreateAbbreviationDTO,
	) error

	UpdateAbbreviation(
		ctx context.Context,
		abbreviation *dto.UpdateAbbreviationDTO,
	) error

	DeleteAbbreviation(
		ctx context.Context,
		standardWord string,
	) error
}
