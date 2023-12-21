package unanswered_question

import (
	"chatbot_be_go/src/application/shared"
	"chatbot_be_go/src/application/unanswered_question/dto"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type IUnansweredQuestionHandler interface {
	ListUnansweredQuestion(c *fiber.Ctx) error
	AnswerQuestion(c *fiber.Ctx) error
}

type unansweredQuestionHandler struct {
	logger                       *logrus.Logger
	unansweredQuestionRepository IUnansweredQuestionRepository
}

func NewUnansweredQuestionHandler(
	logger *logrus.Logger,
	unansweredQuestionRepository IUnansweredQuestionRepository,
) IUnansweredQuestionHandler {
	return &unansweredQuestionHandler{
		logger:                       logger,
		unansweredQuestionRepository: unansweredQuestionRepository,
	}
}

var _ IUnansweredQuestionHandler = &unansweredQuestionHandler{}

// ListUnansweredQuestion godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC05 Menjawab Pertanyaan Manual"
// @Summary Mendapatkan List Pertanyaan Manual
// @Description Mendapatkan List Pertanyaan Manual
// @Param keyword query string false "keyword"
// @Param start query string true "start"
// @Param size query string true "size"
// @Success 200 {object} []dto.UnansweredQuestionDTO
// @Failure 400 {object} string
// @Router /api/unanswered-question [get]
func (u *unansweredQuestionHandler) ListUnansweredQuestion(c *fiber.Ctx) error {
	var counter uint64

	keyword := c.Query("keyword")

	start, _ := strconv.ParseUint(c.Query("start"), 10, 32)
	size, _ := strconv.ParseUint(c.Query("size"), 10, 32)

	if size == 0 {
		size = 10
	}

	callbackSendData := shared.CreateSendCallback(start, &counter, c)

	_, _ = c.Write([]byte("["))

	if totalAll, err := u.unansweredQuestionRepository.ListUnansweredQuestion(
		c.Context(),
		keyword,
		uint(start),
		uint(size),
		callbackSendData,
	); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf(
				"Error on List Admin: %s",
				err.Error(),
			),
		)
	} else {
		_, _ = c.Write([]byte("]"))

		c.Response().Header.Add(
			"Content-Range",
			fmt.Sprintf("items %v-%v/%v", start, counter, totalAll),
		)
	}

	return nil
}

// AnswerQuestion godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC05 Menjawab Pertanyaan Manual"
// @Summary Menjawab Pertanyaan Secara Manual
// @Description Menjawab Pertanyaan Secara Manual
// @Accept json
// @Produce json
// @Param q body dto.AnswerQuestionDTO true "Answer Question"
// @Success 202 {object} dto.ResponseMessageDTO
// @Failure 400 {object} string
// @Router /api/unanswered-question [post]
func (u *unansweredQuestionHandler) AnswerQuestion(c *fiber.Ctx) error {
	userLoggedIn := shared.GetUserData(c)
	if userLoggedIn == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	userQuestionId, err := strconv.ParseUint(c.Params("id", "0"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Format")
	}

	var answerQuestionDTO dto.AnswerQuestionDTO

	if err = json.Unmarshal(c.Body(), &answerQuestionDTO); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Format")
	}

	if err = u.unansweredQuestionRepository.AnswerQuestion(
		c.Context(),
		userQuestionId,
		answerQuestionDTO.Answer,
		userLoggedIn.UserName,
	); err != nil {
		u.logger.WithContext(c.Context()).WithError(err).Error("Error on Answering the Question")
		return c.Status(fiber.StatusInternalServerError).SendString("Something Wrong")
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"message": "User Question Answered Manually Successfully",
	})
}
