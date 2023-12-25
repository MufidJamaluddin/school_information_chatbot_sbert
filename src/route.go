package src

import (
	"chatbot_be_go/src/application"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func RegisterRoute(app *fiber.App, handler *application.Handler, secretResource fiber.Handler) {
	// API Group
	api := app.Group("/api")

	// UC01 Bertanya
	api.Get("/answer", handler.QuestionHandler.GetAnswer)

	webhookWaApi := api.Group("/webhook")
	webhookWaApi.Get("", handler.WhatsAppIntegration.GetWebhook)
	webhookWaApi.Post("", handler.WhatsAppIntegration.PostWebhook)

	// UC02 Admin Login
	loginApi := api.Group("/login")
	loginApi.Post("", handler.LoginHandler.LoginAction)
	loginApi.Delete("", secretResource, handler.LoginHandler.LogoutAction)

	api.Post("/register-admin", secretResource, handler.AdminHandler.RegisterAdmin)
	api.Put("/vector-space-model-reset", secretResource, handler.QuestionHandler.ResetQuestionVector)

	// UC03 Memasukkan Pertanyaan dan Jawaban
	questionHandler := handler.QuestionHandler
	questionApi := api.Group("/question", secretResource)

	questionApi.Put("/:id", questionHandler.UpdateQuestion)
	questionApi.Delete("/:id", questionHandler.DeleteQuestion)
	questionApi.Get("", questionHandler.ListQuestion)
	questionApi.Post("", questionHandler.CreateNewQuestion)

	// UC04 Memasukkan Template Sapaan (Greetings)
	greetingHandler := handler.GreetingHandler
	greetingApi := api.Group("/greeting", secretResource)

	greetingApi.Put("/:id", greetingHandler.UpdateGreeting)
	greetingApi.Delete("/:id", greetingHandler.DeleteGreeting)
	greetingApi.Get("", greetingHandler.ListGreeting)
	greetingApi.Post("", greetingHandler.SaveNewGreeting)

	// UC05 Menjawab Pertanyaan Manual
	unansweredQuestionHandler := handler.UnansweredQuestionHandler
	unansweredQuestionApi := api.Group("/unanswered-question", secretResource)
	unansweredQuestionApi.Post("/:id", unansweredQuestionHandler.AnswerQuestion)
	unansweredQuestionApi.Get("", unansweredQuestionHandler.ListUnansweredQuestion)

	api.Get("/chat-users", secretResource, handler.ChatUserHandler.ListChatUser)

	// UC06 Admin Dashboard
	api.Get("/dashboard-resume", secretResource, handler.DashboardResumeHandler.GetResume)

	// UC07 Manage Admin
	adminHandler := handler.AdminHandler
	adminApi := api.Group("/admin", secretResource)

	adminApi.Get("", adminHandler.ListAdmin)
	adminApi.Put("", adminHandler.UpdateAdmin)

	// UC08 Manage Role Group
	roleGroupHandler := handler.RoleGroupHandler
	roleGroupApi := api.Group("/role-group", secretResource)

	roleGroupApi.Put("/:id", roleGroupHandler.UpdateRoleGroup)
	roleGroupApi.Get("", roleGroupHandler.ListRoleGroup)
	roleGroupApi.Post("", roleGroupHandler.SaveNewRoleGroup)

	// Abbreviation
	abbreviationHandler := handler.AbbreviationHandler
	abbreviationApi := api.Group("/abbreviation", secretResource)
	abbreviationApi.Delete("/:id", abbreviationHandler.DeleteAbbreviation)
	abbreviationApi.Put("", abbreviationHandler.UpdateAbbreviation)
	abbreviationApi.Get("", abbreviationHandler.ListAbbreviation)
	abbreviationApi.Post("", abbreviationHandler.SaveNewAbbreviation)

	// User
	userApi := api.Group("/user")
	userApi.Post("", handler.UserHandler.CreateUser)

	// User Response
	userRespApi := api.Group("/user-response")
	userRespApi.Post("", handler.UserResponseHandler.CreateUserResponse)

	// SPA
	app.Static("", "dist")

	// Swagger Documentation
	app.Get("/docs/*", swagger.HandlerDefault)

	app.Get("/*", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("./dist/index.html")
	})
}
