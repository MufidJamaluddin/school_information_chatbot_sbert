package postgres

import (
	dr "chatbot_be_go/src/application/dashboard_resume"
	"chatbot_be_go/src/application/dashboard_resume/dto"
	"context"
)

type dashboardResumeRepository struct {
	db IDB
}

var _ dr.IDashboardResumeRepository = &dashboardResumeRepository{}

func NewDashboardResumeRepository(db IDB) dr.IDashboardResumeRepository {
	return &dashboardResumeRepository{
		db: db,
	}
}

func (d *dashboardResumeRepository) GetDashboardResume(ctx context.Context) (resume dto.DashboardResumeDTO, err error) {
	sqlDb := d.db.GetSqlDb()

	err = sqlDb.QueryRowContext(
		ctx,
		`WITH totalUser AS (
			SELECT COUNT(*) AS total FROM public."chat_user"
		), totalQuestion AS (
			SELECT COUNT(*) AS total FROM public."question"
		), totalGreeting AS (
			SELECT COUNT(*) AS total FROM public."greeting"
		), totalAdmin AS (
			SELECT COUNT(*) AS total FROM public."admin"
		)
		SELECT
			totalUser.total,
			totalQuestion.total,
			totalGreeting.total,
			totalAdmin.total
		FROM
			totalUser,
			totalQuestion,
			totalGreeting,
			totalAdmin;`,
	).Scan(
		&resume.TotalUser,
		&resume.TotalQuestion,
		&resume.TotalGreeting,
		&resume.TotalAdmin,
	)

	return
}
