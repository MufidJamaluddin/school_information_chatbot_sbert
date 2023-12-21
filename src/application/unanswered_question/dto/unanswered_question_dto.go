package dto

type UnansweredQuestionDTO struct {
	Id       uint64 `json:"id"`
	PhoneNo  string `json:"phoneNo"`
	FullName string `json:"fullName"`
	Question string `json:"question"`
}

type AnswerQuestionDTO struct {
	Answer string `json:"answer"`
}

type ResponseMessageDTO struct {
	Message string `json:"message"`
}
