package question

import "context"

type IQuestionRepository interface {
	FindAnswer(
		nearestAnswer string,
	) (answer string, err error)

	FindAnswerWithSimilarityValue(
		ctx context.Context,
		nearestAnswer string,
	) (
		answer string,
		similarityValue float64,
		err error,
	)

	ListQuestion(
		ctx context.Context,
		keyword string,
		start uint,
		size uint,
		callbackSendData func(questionItem interface{}),
	) (totalAll uint64, err error)

	SaveNewQuestion(
		ctx context.Context,
		question string,
		answer string,
		roleGroupId uint64,
		createdBy string,
	) (uint64, error)

	UpdateQuestion(
		ctx context.Context,
		questionId uint64,
		question string,
		answer string,
		updatedBy string,
	) error

	DeleteQuestion(
		ctx context.Context,
		questionId uint64,
	) error

	TruncateQuestion(
		ctx context.Context,
	) error
}
