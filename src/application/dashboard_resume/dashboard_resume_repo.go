package dashboard_resume

import (
	"chatbot_be_go/src/application/dashboard_resume/dto"
	"context"
)

type IDashboardResumeRepository interface {
	GetDashboardResume(ctx context.Context) (resume dto.DashboardResumeDTO, err error)
}
