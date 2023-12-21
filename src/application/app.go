package application

import (
	abh "chatbot_be_go/src/application/abbreviation"
	"chatbot_be_go/src/application/admin"
	"chatbot_be_go/src/application/chat_user"
	dashboard "chatbot_be_go/src/application/dashboard_resume"
	"chatbot_be_go/src/application/greeting"
	"chatbot_be_go/src/application/login"
	"chatbot_be_go/src/application/question"
	"chatbot_be_go/src/application/role_group"
	"chatbot_be_go/src/application/unanswered_question"
	us "chatbot_be_go/src/application/user"
	usa "chatbot_be_go/src/application/user_response"
	"chatbot_be_go/src/application/whatsapp"
	"chatbot_be_go/src/persistence"
	"chatbot_be_go/src/persistence/config"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	WhatsAppIntegration       whatsapp.IWhatsAppIntegration
	LoginHandler              login.ILoginHandler
	DashboardResumeHandler    dashboard.IDashboardResumeHandler
	UnansweredQuestionHandler unanswered_question.IUnansweredQuestionHandler
	RoleGroupHandler          role_group.IRoleGroupHandler
	ChatUserHandler           chat_user.IChatUserHandler
	AdminHandler              admin.IAdminHandler
	AbbreviationHandler       abh.IAbbreviationHandler
	GreetingHandler           greeting.IGreetingHandler
	QuestionHandler           question.IQuestionHandler
	UserHandler               us.IUserHandler
	UserResponseHandler       usa.IUserResponseHandler
}

func New(
	appConfig *config.AppConfig,
	validator *validator.Validate,
	logger *logrus.Logger,
	httpClient *http.Client,
	persistence *persistence.Persistence,
) *Handler {
	return &Handler{
		WhatsAppIntegration: whatsapp.New(
			appConfig,
			validator,
			logger,
			httpClient,
		),
		LoginHandler: login.NewLoginHandler(
			logger,
			appConfig,
			persistence.LoginRepository,
		),
		DashboardResumeHandler: dashboard.NewDashboardResumeHandler(
			logger,
			persistence.DashboardResumeRepository,
		),
		UnansweredQuestionHandler: unanswered_question.NewUnansweredQuestionHandler(
			logger,
			persistence.UnansweredQuestionRepository,
		),
		RoleGroupHandler: role_group.NewRoleGroupHandler(
			logger,
			persistence.RoleGroupRepository,
		),
		ChatUserHandler: chat_user.NewChatUserHandler(
			logger,
			persistence.ChatUserRepository,
		),
		AdminHandler: admin.NewAdminHandler(
			logger,
			persistence.AdminRepository,
		),
		AbbreviationHandler: abh.NewAbbreviationHandler(
			logger,
			persistence.AbbreviationRepository,
		),
		GreetingHandler: greeting.NewGreetingHandler(
			logger,
			persistence.GreetingRepository,
		),
		QuestionHandler: question.NewQuestionHandler(
			logger,
			persistence.QuestionRepository,
			persistence.GreetingRepository,
		),
		UserHandler: us.NewUserHandler(
			logger,
			persistence.UserRepository,
		),
		UserResponseHandler: usa.NewUserResponseHandler(
			logger,
			persistence.UserResponseRepository,
		),
	}
}
