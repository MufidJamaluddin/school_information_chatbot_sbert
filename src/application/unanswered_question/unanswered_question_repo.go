package unanswered_question

import "context"

type IUnansweredQuestionRepository interface {
	ListUnansweredQuestion(
		ctx context.Context,
		keyword string,
		start uint,
		size uint,
		callbackSendData func(adminItem interface{}),
	) (totalAll uint64, err error)

	AnswerQuestion(
		ctx context.Context,
		userQuestionId uint64,
		answer string,
		answerByUserName string,
	) error
}
