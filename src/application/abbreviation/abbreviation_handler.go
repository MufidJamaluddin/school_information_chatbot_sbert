package abbreviation

import (
	dto "chatbot_be_go/src/application/abbreviation/dto"
	"chatbot_be_go/src/application/shared"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type IAbbreviationHandler interface {
	ListAbbreviation(c *fiber.Ctx) error
	SaveNewAbbreviation(c *fiber.Ctx) error
	UpdateAbbreviation(c *fiber.Ctx) error
	DeleteAbbreviation(c *fiber.Ctx) error
}

type abbreviationHandler struct {
	logger           *logrus.Logger
	abbreviationRepo IAbbreviationRepository
}

var _ IAbbreviationHandler = &abbreviationHandler{}

func NewAbbreviationHandler(
	logger *logrus.Logger,
	abbreviationRepo IAbbreviationRepository,
) IAbbreviationHandler {
	return &abbreviationHandler{
		logger:           logger,
		abbreviationRepo: abbreviationRepo,
	}
}

// ListAbbreviation godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC09 Manage Abbreviation"
// @Summary Get the List of Abbreviation's Data
// @Description Show List of Abbreviation's Data
// @Param keyword query string false "keyword"
// @Param start query string true "start"
// @Param size query string true "size"
// @Success 200 {object} []dto.AbbreviationItemDTO
// @Failure 400 {object} string
// @Router /api/abbreviation [get]
func (s *abbreviationHandler) ListAbbreviation(c *fiber.Ctx) error {
	var counter uint64

	keyword := c.Query("keyword")

	start, _ := strconv.ParseUint(c.Query("start"), 10, 32)
	size, _ := strconv.ParseUint(c.Query("size"), 10, 32)

	if size == 0 {
		size = 10
	}

	callbackSendData := shared.CreateSendCallback(start, &counter, c)

	_, _ = c.Write([]byte("["))

	if totalAll, err := s.abbreviationRepo.ListAbbreviation(
		c.Context(),
		keyword,
		uint(start),
		uint(size),
		callbackSendData,
	); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf(
				"Error on Abbreviationn: %s",
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

// SaveNewAbbreviation godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC09 Manage Abbreviation"
// @Summary Save New Abbreviation
// @Description Save New Abbreviation Data
// @Accept json
// @Produce json
// @Param q body dto.CreateAbbreviationDTO true "New Abbreviation Data"
// @Success 202 {object} dto.ResponseMessageDTO
// @Failure 400 {object} string
// @Router /api/abbreviation [post]
func (s *abbreviationHandler) SaveNewAbbreviation(c *fiber.Ctx) error {
	var abbreviationDtoData dto.CreateAbbreviationDTO

	userLoggedIn := shared.GetUserData(c)
	if userLoggedIn == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	if err := json.Unmarshal(c.Body(), &abbreviationDtoData); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Format")
	}

	abbreviationDtoData.CreatedBy = userLoggedIn.UserName

	if err := s.abbreviationRepo.SaveNewAbbreviation(c.Context(), &abbreviationDtoData); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf(
				"Error on Save Abbreviation: %s",
				err.Error(),
			),
		)
	} else {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id":      abbreviationDtoData.StandardWord,
			"message": "Abbreviation Successfully Created!",
		})
	}
}

// UpdateAbbreviation godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC09 Manage Abbreviation"
// @Summary Update Abbreviation Data
// @Description Update Abbreviation Data
// @Accept json
// @Produce json
// @Param q body dto.UpdateAbbreviationDTO true "Abbreviation Data"
// @Success 202 {object} dto.ResponseMessageDTO
// @Failure 400 {object} string
// @Router /api/abbreviation [put]
func (s *abbreviationHandler) UpdateAbbreviation(c *fiber.Ctx) error {
	var updateAbbreviationItemData dto.UpdateAbbreviationDTO

	userLoggedIn := shared.GetUserData(c)
	if userLoggedIn == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	if err := json.Unmarshal(c.Body(), &updateAbbreviationItemData); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Format")
	}

	updateAbbreviationItemData.UpdatedBy = userLoggedIn.UserName

	if err := s.abbreviationRepo.UpdateAbbreviation(c.Context(), &updateAbbreviationItemData); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf(
				"Error on Save Greeting: %s",
				err.Error(),
			),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Abbreviation Successfully Updated!",
	})
}

// DeleteAbbreviation godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC09 Manage Abbreviation"
// @Summary Delete one Abbreviation by id
// @Description Delete one Abbreviation by id
// @Param id path string true "Abbreviation Standard Word"
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.ResponseMessageDTO
// @Failure 400 {object} string
// @Failure 404 {object} string
// @Router /api/abbreviation/{id} [delete]
func (s *abbreviationHandler) DeleteAbbreviation(c *fiber.Ctx) error {
	standardWord := c.Params("id", "0")

	if err := s.abbreviationRepo.DeleteAbbreviation(c.Context(), standardWord); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf(
				"Error on Save Abbreviation: %s",
				err.Error(),
			),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":      standardWord,
		"message": "Abbreviation Successfully Deleted!",
	})
}
