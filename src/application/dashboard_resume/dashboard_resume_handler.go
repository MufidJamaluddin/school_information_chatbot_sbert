package dashboard_resume

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type IDashboardResumeHandler interface {
	GetResume(c *fiber.Ctx) error
}

type dashboardResumeHandler struct {
	logger                    *logrus.Logger
	dashboardResumeRepository IDashboardResumeRepository
}

var _ IDashboardResumeHandler = &dashboardResumeHandler{}

func NewDashboardResumeHandler(
	logger *logrus.Logger,
	dashboardResumeRepository IDashboardResumeRepository,
) IDashboardResumeHandler {
	return &dashboardResumeHandler{
		logger:                    logger,
		dashboardResumeRepository: dashboardResumeRepository,
	}
}

// GetResume godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC06 Admin Dashboard"
// @Summary Dashboard Resume Data
// @Description Get Dashboard Resume Data
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.DashboardResumeDTO
// @Failure 400 {object} string
// @Router /api/dashboard-resume [get]
func (d *dashboardResumeHandler) GetResume(c *fiber.Ctx) error {
	if resume, err := d.dashboardResumeRepository.GetDashboardResume(c.Context()); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Unauthorized")
	} else {
		return c.Status(fiber.StatusOK).JSON(&resume)
	}
}
