package dto

import "time"

type CreateAbbreviationDTO struct {
	StandardWord         string   `json:"standardWord"`
	ListAbbreviationTerm []string `json:"listAbbreviationTerm"`
	CreatedBy            string   `json:"-"`
}

type UpdateAbbreviationDTO struct {
	StandardWord         string   `json:"standardWord"`
	ListAbbreviationTerm []string `json:"listAbbreviationTerm"`
	UpdatedBy            string   `json:"-"`
}

type AbbreviationItemDTO struct {
	StandardWord         string     `json:"standardWord"`
	ListAbbreviationTerm []string   `json:"listAbbreviationTerm"`
	CreatedAt            time.Time  `json:"createdAt"`
	UpdatedAt            *time.Time `json:"updatedAt"`
	CreatedBy            string     `json:"createdBy"`
	UpdatedBy            string     `json:"updatedBy"`
}

type ResponseMessageDTO struct {
	Id      uint64 `json:"id"`
	Message string `json:"message"`
}
