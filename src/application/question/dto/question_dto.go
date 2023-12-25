package dto

type QuestionItemDTO struct {
	QuestionId  uint64 `json:"id"`
	Question    string `json:"question"`
	Answer      string `json:"answer"`
	RoleGroupId uint64 `json:"roleGroupId"`
	CreatedBy   string `json:"createdBy"`
	UpdatedBy   string `json:"updatedBy"`
}

type QuestionItemDTOWithVector struct {
	QuestionId uint64    `json:"id"`
	Question   string    `json:"question"`
	Vector     []float64 `json:"vector"`
}

type CreateQuestionDTO struct {
	Question    string `json:"question" validate:"nonzero,max=255"`
	Answer      string `json:"answer" validate:"nonzero,max=2550"`
	RoleGroupId uint64 `json:"-"`
	CreatedBy   string `json:"-"`
}

type CreateQuestionResponseDTO struct {
	QuestionId uint64 `json:"questionId"`
	Message    string `json:"message"`
}

type UpdateQuestionDTO struct {
	Question  string `json:"question" validate:"nonzero,max=255"`
	Answer    string `json:"answer" validate:"nonzero,max=2550"`
	UpdatedBy string `json:"-"`
}

type UpdateQuestionResponseDTO struct {
	Message string `json:"message"`
}

type AnswerQuestionDTO struct {
	Answer  string `json:"answer"`
	Message string `json:"message,omitempty"`
}
