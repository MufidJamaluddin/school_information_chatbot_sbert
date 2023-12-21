package greeting

import (
	"chatbot_be_go/src/application/greeting/dto"
	"chatbot_be_go/src/application/shared"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type IGreetingHandler interface {
	ListGreeting(c *fiber.Ctx) error
	SaveNewGreeting(c *fiber.Ctx) error
	UpdateGreeting(c *fiber.Ctx) error
	DeleteGreeting(c *fiber.Ctx) error
}

type greetingHandler struct {
	logger             *logrus.Logger
	greetingRepository IGreetingRepository
}

func NewGreetingHandler(
	logger *logrus.Logger,
	greetingRepository IGreetingRepository,
) IGreetingHandler {
	return &greetingHandler{
		logger:             logger,
		greetingRepository: greetingRepository,
	}
}

var _ IGreetingHandler = &greetingHandler{}

// ListGreeting godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC04 Memasukkan Template Sapaan (Greetings)"
// @Summary Get the List Greeting Data
// @Description Show List Greeting Data
// @Param keyword query string false "keyword"
// @Param start query string true "start"
// @Param size query string true "size"
// @Success 200 {object} []dto.GreetingItemDTO
// @Failure 400 {object} string
// @Router /api/greeting [get]
func (g *greetingHandler) ListGreeting(c *fiber.Ctx) error {
	var counter uint64

	keyword := c.Query("keyword")

	start, _ := strconv.ParseUint(c.Query("start"), 10, 32)
	size, _ := strconv.ParseUint(c.Query("size"), 10, 32)

	if size == 0 {
		size = 10
	}

	callbackSendData := shared.CreateSendCallback(start, &counter, c)

	_, _ = c.Write([]byte("["))

	if totalAll, err := g.greetingRepository.ListGreeting(
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

// SaveNewGreeting godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC04 Memasukkan Template Sapaan (Greetings)"
// @Summary Save New Greeting
// @Description Save New Greeting Data
// @Accept json
// @Produce json
// @Param q body dto.CreateGreetingDTO true "New Greeting"
// @Success 202 {object} dto.ResponseMessageDTO
// @Failure 400 {object} string
// @Router /api/greeting [post]
func (g *greetingHandler) SaveNewGreeting(c *fiber.Ctx) error {
	var greeting dto.CreateGreetingDTO

	userLoggedIn := shared.GetUserData(c)
	if userLoggedIn == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	if err := json.Unmarshal(c.Body(), &greeting); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Format")
	}

	greeting.CreatedBy = userLoggedIn.UserName

	if newId, err := g.greetingRepository.SaveNewGreeting(c.Context(), &greeting); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf(
				"Error on Save Greeting: %s",
				err.Error(),
			),
		)
	} else {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id":      newId,
			"message": "Greeting Successfully Created!",
		})
	}
}

// UpdateGreeting godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC04 Memasukkan Template Sapaan (Greetings)"
// @Summary Update Greeting Data
// @Description Update Greeting Data
// @Accept json
// @Produce json
// @param id path int true "Greeting ID"
// @Param q body dto.UpdateGreetingDTO true "Greeting Data"
// @Success 202 {object} dto.ResponseMessageDTO
// @Failure 400 {object} string
// @Router /api/greeting/{id} [put]
func (g *greetingHandler) UpdateGreeting(c *fiber.Ctx) error {
	var greeting dto.UpdateGreetingDTO

	userLoggedIn := shared.GetUserData(c)
	if userLoggedIn == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	greetingId, err := strconv.ParseUint(c.Params("id", "0"), 10, 64)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(c.Body(), &greeting); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Format")
	}

	greeting.UpdateBy = userLoggedIn.UserName

	if err = g.greetingRepository.UpdateGreeting(c.Context(), greetingId, &greeting); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf(
				"Error on Save Greeting: %s",
				err.Error(),
			),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":      greetingId,
		"message": "Greeting Successfully Updated!",
	})
}

// DeleteGreeting godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC04 Memasukkan Template Sapaan (Greetings)"
// @Summary Delete one Greeting by id
// @Description Delete one Greeting by id
// @Param id path int true "Greeting ID"
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.ResponseMessageDTO
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Router /api/greeting/{id} [delete]
func (g *greetingHandler) DeleteGreeting(c *fiber.Ctx) error {
	greetingId, err := strconv.ParseUint(c.Params("id", "0"), 10, 64)
	if err != nil {
		return err
	}

	if err = g.greetingRepository.DeleteGreeting(c.Context(), greetingId); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf(
				"Error on Save Greeting: %s",
				err.Error(),
			),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":      greetingId,
		"message": "Greeting Successfully Deleted!",
	})
}
