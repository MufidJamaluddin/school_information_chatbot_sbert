package role_group

import (
	"chatbot_be_go/src/application/role_group/dto"
	"chatbot_be_go/src/application/shared"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type IRoleGroupHandler interface {
	ListRoleGroup(c *fiber.Ctx) error
	SaveNewRoleGroup(c *fiber.Ctx) error
	UpdateRoleGroup(c *fiber.Ctx) error
}

type roleGroupHandler struct {
	logger              *logrus.Logger
	roleGroupRepository IRoleGroupRepository
}

func NewRoleGroupHandler(
	logger *logrus.Logger,
	roleGroupRepository IRoleGroupRepository,
) IRoleGroupHandler {
	return &roleGroupHandler{
		logger:              logger,
		roleGroupRepository: roleGroupRepository,
	}
}

var _ IRoleGroupHandler = &roleGroupHandler{}

// ListRoleGroup godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC08 Manage Role Group"
// @Summary Get the List Role Group Data
// @Description Show List Role Group Data
// @Param keyword query string false "keyword"
// @Param start query string true "start"
// @Param size query string true "size"
// @Success 200 {object} []dto.RoleItemGroupDTO
// @Failure 400 {object} string
// @Router /api/role-group [get]
func (r *roleGroupHandler) ListRoleGroup(c *fiber.Ctx) error {
	var counter uint64

	keyword := c.Query("keyword")

	start, _ := strconv.ParseUint(c.Query("start"), 10, 32)
	size, _ := strconv.ParseUint(c.Query("size"), 10, 32)

	if size == 0 {
		size = 10
	}

	_, _ = c.Write([]byte("["))

	callbackSendData := shared.CreateSendCallback(start, &counter, c)

	if totalAll, err := r.roleGroupRepository.ListRoleGroup(
		c.Context(),
		keyword,
		uint(start),
		uint(size),
		callbackSendData,
	); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf(
				"Error on List Role Group: %s",
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

// SaveNewRoleGroup godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC08 Manage Role Group"
// @Summary Save New Role Group
// @Description Save New Role Group Data
// @Accept json
// @Produce json
// @Param q body dto.CreateRoleGroupDTO true "New Role Group"
// @Success 202 {object} dto.ResponseMessageDTO
// @Failure 400 {object} string
// @Router /api/role-group [post]
func (r *roleGroupHandler) SaveNewRoleGroup(c *fiber.Ctx) error {
	var roleGroup dto.CreateRoleGroupDTO

	userLoggedIn := shared.GetUserData(c)
	if userLoggedIn == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	roleGroup.CreatedBy = userLoggedIn.UserName

	if err := json.Unmarshal(c.Body(), &roleGroup); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Format")
	}

	if newId, err := r.roleGroupRepository.SaveNewRoleGroup(c.Context(), &roleGroup); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf(
				"Error on Save Role Group: %s",
				err.Error(),
			),
		)
	} else {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"id":      newId,
			"message": "Role Group Successfully Created!",
		})
	}
}

// UpdateRoleGroup godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC08 Manage Role Group"
// @Summary Update Role Group Data
// @Description Update Role Group Data
// @Accept json
// @Produce json
// @param id path int true "Role Group ID"
// @Param q body dto.UpdateRoleGroupDTO true "Role Group Data"
// @Success 202 {object} dto.ResponseMessageDTO
// @Failure 400 {object} string
// @Router /api/role-group/{id} [put]
func (r *roleGroupHandler) UpdateRoleGroup(c *fiber.Ctx) error {
	var roleGroup dto.UpdateRoleGroupDTO

	userLoggedIn := shared.GetUserData(c)
	if userLoggedIn == nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
	}

	roleGroupId, err := strconv.ParseUint(c.Params("id", "0"), 10, 64)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(c.Body(), &roleGroup); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Format")
	}

	roleGroup.UpdatedBy = userLoggedIn.UserName

	if err = r.roleGroupRepository.UpdateRoleGroup(c.Context(), roleGroupId, &roleGroup); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(
			fmt.Sprintf(
				"Error on Save Role Group: %s",
				err.Error(),
			),
		)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":      roleGroupId,
		"message": "Role Group Successfully Updated!",
	})
}
