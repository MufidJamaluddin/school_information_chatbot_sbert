package admin

import (
	"chatbot_be_go/src/application/user_response/dto"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/validator.v2"
)

type IUserResponseHandler interface {
	CreateUserResponse(c *fiber.Ctx) error
}

type userAnswerHandler struct {
	logger         *logrus.Logger
	userAnswerRepo IUserResponseRepository
}

func NewUserResponseHandler(
	logger *logrus.Logger,
	userAnswerRepo IUserResponseRepository,
) IUserResponseHandler {
	return &userAnswerHandler{
		logger:         logger,
		userAnswerRepo: userAnswerRepo,
	}
}

var _ IUserResponseHandler = &userAnswerHandler{}

// CreateUser godoc
// @Tags "CreateUserResponse"
// @Summary Create User Response
// @Description Save New User Response
// @Accept json
// @Produce json
// @Param q body dto.UserResponse true "User Response Data"
// @Success 202 {object} dto.User
// @Failure 400 {object} string
// @Router /api/user-response [post]
func (uh *userAnswerHandler) CreateUserResponse(c *fiber.Ctx) (err error) {
	var userAnswerDto dto.UserResponse

	if err = json.Unmarshal(c.Body(), &userAnswerDto); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Format")
	}

	if err = validator.Validate(&userAnswerDto); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(
			fmt.Sprintf("Invalid Format: %v", err),
		)
	}

	if err = uh.userAnswerRepo.CreateNewUserResponse(c.Context(), &userAnswerDto); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf(
				"Error on Save User Response: %s",
				err.Error(),
			),
		)
	} else {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User Answer Succesfully Recorded!",
		})
	}
}
