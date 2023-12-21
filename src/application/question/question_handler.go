package question

import (
	"chatbot_be_go/src/application/greeting"
	"chatbot_be_go/src/application/question/dto"
	"chatbot_be_go/src/application/shared"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/validator.v2"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type IQuestionHandler interface {
	GetAnswer(c *fiber.Ctx) error
	ListQuestion(c *fiber.Ctx) error
	CreateNewQuestion(c *fiber.Ctx) error
	UpdateQuestion(c *fiber.Ctx) error
	DeleteQuestion(c *fiber.Ctx) error
}

type questionHandler struct {
	logger             *logrus.Logger
	questionRepository IQuestionRepository
	greetingRepository greeting.IGreetingRepository
}

var _ IQuestionHandler = &questionHandler{}

func NewQuestionHandler(
	logger *logrus.Logger,
	questionRepository IQuestionRepository,
	greetingRepository greeting.IGreetingRepository,
) IQuestionHandler {
	return &questionHandler{
		logger:             logger,
		questionRepository: questionRepository,
		greetingRepository: greetingRepository,
	}
}

// GetAnswer godoc
// @Tags "UC01 Bertanya"
// @Summary Answer the Incoming Question
// @Description Answer the Incoming Question
// @Param question query string true "Question"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /api/answer [get]
func (q *questionHandler) GetAnswer(c *fiber.Ctx) (err error) {
	var answerChat string
	var answerGreeting string

	var answer dto.AnswerQuestionDTO
	question := c.Query("question")
	name := c.Query("name")

	if answerGreeting, err = q.greetingRepository.FindCurrentGreeting(c.Context()); err != nil {
		answerGreeting = ""
	}

	if answerChat, err = q.questionRepository.FindAnswer(question); err != nil {
		q.logger.Error(err)

		return c.Status(fiber.StatusNotFound).SendString(
			fmt.Sprintf(
				"Error on Save Question: %s",
				err.Error(),
			),
		)
	}

	if answerChat == "" {
		answerChat = "Jawaban tidak ketemu"
	}

	if answerGreeting == "" {
		answerGreeting = "Baik kak {name}"
	}

	answerGreeting = strings.TrimSpace(strings.ReplaceAll(answerGreeting, "{name}", name))

	answer.Answer = fmt.Sprintf("%s. %s", answerGreeting, answerChat)

	return c.JSON(&answer)
}

// ListQuestion godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC03 Memasukkan Pertanyaan dan Jawaban"
// @Summary Get the List Greeting Data
// @Description Show List Greeting Data
// @Param keyword query string false "keyword"
// @Param start query string true "start"
// @Param size query string true "size"
// @Success 200 {object} []dto.QuestionItemDTO
// @Failure 400 {object} string
// @Router /api/question [get]
func (q *questionHandler) ListQuestion(c *fiber.Ctx) error {
	var counter uint64

	keyword := c.Query("keyword")

	start, _ := strconv.ParseUint(c.Query("start"), 10, 32)
	size, _ := strconv.ParseUint(c.Query("size"), 10, 32)

	if size == 0 {
		size = 10
	}

	_, _ = c.Write([]byte("["))

	callbackSendData := shared.CreateSendCallback(start, &counter, c)

	if totalAll, err := q.questionRepository.ListQuestion(
		c.Context(),
		keyword,
		uint(start),
		uint(size),
		callbackSendData,
	); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf(
				"Error on List Question: %s",
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

// CreateNewQuestion godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC03 Memasukkan Pertanyaan dan Jawaban"
// @Summary Save New Question
// @Description Save New Question to Answer the Repetitive Question
// @Accept json
// @Produce json
// @Param q body dto.CreateQuestionDTO true "New Question"
// @Success 202 {object} dto.CreateQuestionResponseDTO
// @Failure 400 {object} string
// @Router /api/question [post]
func (q *questionHandler) CreateNewQuestion(c *fiber.Ctx) (err error) {
	var createDataDto dto.CreateQuestionDTO
	var responseDataDto dto.CreateQuestionResponseDTO

	userLoggedIn := shared.GetUserData(c)
	if userLoggedIn == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	if err = json.Unmarshal(c.Body(), &createDataDto); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Format")
	}

	if err = validator.Validate(&createDataDto); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(
			fmt.Sprintf("Invalid Format: %v", err),
		)
	}

	createDataDto.CreatedBy = userLoggedIn.UserName
	createDataDto.RoleGroupId = userLoggedIn.RoleGroupId

	if responseDataDto.QuestionId, err = q.questionRepository.SaveNewQuestion(
		c.Context(),
		createDataDto.Question,
		createDataDto.Answer,
		createDataDto.RoleGroupId,
		createDataDto.CreatedBy,
	); err != nil {
		q.logger.Error(err)

		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf(
				"Error on Save Question: %s",
				err.Error(),
			),
		)
	}

	responseDataDto.Message = "The Question Saved Successfully"

	return c.Status(fiber.StatusCreated).JSON(&responseDataDto)
}

// UpdateQuestion godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC03 Memasukkan Pertanyaan dan Jawaban"
// @Summary Update Question
// @Description Change Question and Answer
// @Accept json
// @Produce json
// @Param id path string true "Question ID"
// @Param q body dto.UpdateQuestionDTO true "New Question"
// @Success 202 {object} dto.UpdateQuestionResponseDTO
// @Failure 400 {object} string
// @Router /api/question/{id} [put]
func (q *questionHandler) UpdateQuestion(c *fiber.Ctx) error {
	userLoggedIn := shared.GetUserData(c)
	if userLoggedIn == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	questionId, err := strconv.ParseUint(c.Params("id", "0"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Format")
	}

	var updateDataDto dto.UpdateQuestionDTO
	var responseUpdateDataDto dto.UpdateQuestionResponseDTO

	if err = json.Unmarshal(c.Body(), &updateDataDto); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Format")
	}

	if err = validator.Validate(&updateDataDto); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(
			fmt.Sprintf("Invalid Format: %v", err),
		)
	}

	updateDataDto.UpdatedBy = userLoggedIn.UserName

	if err = q.questionRepository.UpdateQuestion(
		c.Context(),
		questionId,
		updateDataDto.Question,
		updateDataDto.Answer,
		updateDataDto.UpdatedBy,
	); err != nil {
		q.logger.Error(err)

		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf(
				"Error on Update Question: %s",
				err.Error(),
			),
		)
	}

	responseUpdateDataDto.Message = "Question Updated Successfully"

	return c.Status(fiber.StatusOK).JSON(&responseUpdateDataDto)
}

// DeleteQuestion godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC03 Memasukkan Pertanyaan dan Jawaban"
// @Summary Delete one Question by id
// @Description Delete one Question by id
// @Param id path int true "Question ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.UpdateQuestionResponseDTO
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Router /api/question/{id} [delete]
func (q *questionHandler) DeleteQuestion(c *fiber.Ctx) error {
	questionId, err := strconv.ParseUint(c.Params("id", "0"), 10, 64)
	if err != nil {
		return err
	}

	if err = q.questionRepository.DeleteQuestion(c.Context(), questionId); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf(
				"Error on Save Greeting: %s",
				err.Error(),
			),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":      questionId,
		"message": "Question Successfully Deleted!",
	})
}
