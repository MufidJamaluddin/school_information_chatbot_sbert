package admin

import (
	"chatbot_be_go/src/application/admin/dto"
	"chatbot_be_go/src/application/shared"
	"encoding/json"
	"fmt"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/validator.v2"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type IAdminHandler interface {
	RegisterAdmin(c *fiber.Ctx) error
	ListAdmin(c *fiber.Ctx) error
	UpdateAdmin(c *fiber.Ctx) error
}

type adminHandler struct {
	logger          *logrus.Logger
	adminRepository IAdminRepository
}

func NewAdminHandler(
	logger *logrus.Logger,
	adminRepository IAdminRepository,
) IAdminHandler {
	return &adminHandler{
		logger:          logger,
		adminRepository: adminRepository,
	}
}

var _ IAdminHandler = &adminHandler{}

// RegisterAdmin godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC02 Admin Login"
// @Summary Save New Admin
// @Description Save New Admin Data
// @Accept json
// @Produce json
// @Param q body dto.Admin true "New Admin"
// @Success 202 {object} dto.ResponseMessage
// @Failure 400 {object} string
// @Router /api/register-admin [post]
func (a *adminHandler) RegisterAdmin(c *fiber.Ctx) (err error) {
	var admin dto.Admin

	userLoggedIn := shared.GetUserData(c)
	if userLoggedIn == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	if err = json.Unmarshal(c.Body(), &admin); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Format")
	}

	if err = validator.Validate(&admin); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(
			fmt.Sprintf("Invalid Format: %v", err),
		)
	}

	if admin.HashedPassword, err = bcrypt.GenerateFromPassword([]byte(admin.Password), 8); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Format")
	}

	if err = a.adminRepository.CreateNewAdmin(c.Context(), admin, userLoggedIn.UserName); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf(
				"Error on Save Admin: %s",
				err.Error(),
			),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Admin Successfully Created!",
	})
}

// ListAdmin godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC07 Manage Admin"
// @Summary Get the List Admin Data
// @Description Show List Admin Data
// @Param keyword query string false "keyword"
// @Param start query string true "start"
// @Param size query string true "size"
// @Success 200 {object} []dto.ListAdminItem
// @Failure 400 {object} string
// @Router /api/admin [get]
func (a *adminHandler) ListAdmin(c *fiber.Ctx) error {
	var counter uint64

	keyword := c.Query("keyword")

	start, _ := strconv.ParseUint(c.Query("start"), 10, 32)
	size, _ := strconv.ParseUint(c.Query("size"), 10, 32)

	if size == 0 {
		size = 10
	}

	callbackSendData := shared.CreateSendCallback(start, &counter, c)

	_, _ = c.Write([]byte("["))

	if totalAll, err := a.adminRepository.ListAdmin(
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

// UpdateAdmin godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC07 Manage Admin"
// @Summary Update Admin Data (Profile)
// @Description Update Admin Data (Profile)
// @Accept json
// @Produce json
// @Param q body dto.Admin true "Admin Data"
// @Success 202 {object} dto.ResponseMessage
// @Failure 400 {object} string
// @Router /api/admin [put]
func (a *adminHandler) UpdateAdmin(c *fiber.Ctx) (err error) {
	var admin dto.Admin

	userLoggedIn := shared.GetUserData(c)
	if userLoggedIn == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	if err = json.Unmarshal(c.Body(), &admin); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Format")
	}

	if err = validator.Validate(&admin); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString(
			fmt.Sprintf("Invalid Format: %v", err),
		)
	}

	if admin.HashedPassword, err = bcrypt.GenerateFromPassword([]byte(admin.Password), 8); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Format")
	}

	if err := a.adminRepository.UpdateAdmin(c.Context(), admin.Username, admin, userLoggedIn.UserName); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf(
				"Error on Save Admin: %s",
				err.Error(),
			),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Admin Successfully Updated!",
	})
}
