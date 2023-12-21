package persistence

import (
	abh "chatbot_be_go/src/application/abbreviation"
	ads "chatbot_be_go/src/application/admin"
	cu "chatbot_be_go/src/application/chat_user"
	dr "chatbot_be_go/src/application/dashboard_resume"
	gr "chatbot_be_go/src/application/greeting"
	lg "chatbot_be_go/src/application/login"
	q "chatbot_be_go/src/application/question"
	rg "chatbot_be_go/src/application/role_group"
	uq "chatbot_be_go/src/application/unanswered_question"
	us "chatbot_be_go/src/application/user"
	usa "chatbot_be_go/src/application/user_response"
	dm "chatbot_be_go/src/domain"
	"chatbot_be_go/src/persistence/config"
	"chatbot_be_go/src/persistence/postgres"

	"github.com/sirupsen/logrus"
)

type Persistence struct {
	LoginRepository              lg.ILoginRepository
	DashboardResumeRepository    dr.IDashboardResumeRepository
	UnansweredQuestionRepository uq.IUnansweredQuestionRepository
	RoleGroupRepository          rg.IRoleGroupRepository
	ChatUserRepository           cu.IChatUserRepository
	AdminRepository              ads.IAdminRepository
	AbbreviationRepository       abh.IAbbreviationRepository
	GreetingRepository           gr.IGreetingRepository
	QuestionRepository           q.IQuestionRepository
	UserResponseRepository       usa.IUserResponseRepository
	UserRepository               us.IUserRepository
}

func New(
	vectorizer dm.ISBertVectorizer,
	sqlDbConf config.SqlDbConf,
	logger *logrus.Logger,
) *Persistence {
	db := postgres.New(
		sqlDbConf,
		logger,
	)

	loginRepository := postgres.NewLoginRepository(db)
	dashboardResumeRepository := postgres.NewDashboardResumeRepository(db)
	unansweredQuestionRepository := postgres.NewUnansweredQuestionRepository(db)
	roleGroupRepository := postgres.NewRoleGroupRepository(db)
	chatUserRepository := postgres.NewChatUserRepository(db)
	adminRepository := postgres.NewAdminRepository(db)
	abbreviationRepository := postgres.NewAbbreviationRepository(db)
	greetingRepository := postgres.NewGreetingRepository(db)
	questionRepository := postgres.NewQuestionRepository(db, vectorizer)
	userRepository := postgres.NewUserRepository(db)
	userResponseRepository := postgres.NewUserResponseRepository(db)

	return &Persistence{
		LoginRepository:              loginRepository,
		DashboardResumeRepository:    dashboardResumeRepository,
		UnansweredQuestionRepository: unansweredQuestionRepository,
		RoleGroupRepository:          roleGroupRepository,
		ChatUserRepository:           chatUserRepository,
		AdminRepository:              adminRepository,
		AbbreviationRepository:       abbreviationRepository,
		GreetingRepository:           greetingRepository,
		QuestionRepository:           questionRepository,
		UserRepository:               userRepository,
		UserResponseRepository:       userResponseRepository,
	}
}
