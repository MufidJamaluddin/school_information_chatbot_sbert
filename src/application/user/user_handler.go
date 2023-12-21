package admin

import (
	"chatbot_be_go/src/application/user/dto"
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gopkg.in/validator.v2"
)

type IUserHandler interface {
	CreateUser(c *fiber.Ctx) error
}

type userHandler struct {
	logger   *logrus.Logger
	userRepo IUserRepository
}

func NewUserHandler(
	logger *logrus.Logger,
	userRepo IUserRepository,
) IUserHandler {
	return &userHandler{
		logger:   logger,
		userRepo: userRepo,
	}
}

var _ IUserHandler = &userHandler{}

// CreateUser godoc
// @Tags "CreateUser"
// @Summary Create User
// @Description Create User & Get ID
// @Accept json
// @Produce json
// @Param q body dto.User true "User Data"
// @Success 202 {object} dto.User
// @Failure 400 {object} string
// @Router /api/user [post]
func (uh *userHandler) CreateUser(c *fiber.Ctx) (err error) {
	var userDto dto.User

	if err = json.Unmarshal(c.Body(), &userDto); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Format")
	}

	if err = validator.Validate(&userDto); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(
			fmt.Sprintf("Invalid Format: %v", err),
		)
	}

	if err = uh.userRepo.CreateNewUser(c.Context(), &userDto); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf(
				"Error on Save User: %s",
				err.Error(),
			),
		)
	} else {
		return c.Status(fiber.StatusCreated).JSON(&userDto)
	}
}
