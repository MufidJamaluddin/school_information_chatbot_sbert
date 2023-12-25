package main

import (
	_ "chatbot_be_go/docs"
	"chatbot_be_go/src"
	"chatbot_be_go/src/application"
	dm "chatbot_be_go/src/domain"
	"chatbot_be_go/src/persistence"
	appConf "chatbot_be_go/src/persistence/config"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	midLog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// @title Chat Bot Back End API
// @version 1.0
// @description Chat Bot Back End API
// @BasePath /
func main() {
	_ = godotenv.Load()

	app := fiber.New()
	logger := logrus.New()

	config := appConf.New()
	appValidator := validator.New()

	httpClient := &http.Client{
		Timeout: config.Http.S2STimeoutSecs,
	}

	var vectorizer dm.ISBertVectorizer

	if config.IsFineTuneModel {
		vectorizer = dm.NewSBertVectorizer("/fine_tuned_model", "model.pkl")
	} else {
		vectorizer = dm.NewSBertVectorizer("/ori_model", "model.pkl")
	}

	persistenceObj := persistence.New(
		vectorizer,
		config.SqlDb,
		logger,
	)

	handler := application.New(
		config,
		appValidator,
		logger,
		httpClient,
		persistenceObj,
	)

	// Fix Anomaly Different Result per Build
	if err := persistenceObj.QuestionRepository.ResetSBERTVectorQuestion(
		context.Background(),
	); err != nil {
		log.Fatal(err)
	}

	jwtMiddleware := jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: config.SecretKey,
		},
		TokenLookup: fmt.Sprintf(
			"header:%s,cookie:%s",
			config.HeaderTokenKeyName,
			config.CookieKeyName,
		),
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
		},
	})

	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))

	if config.IsLogged {
		app.Use(midLog.New(midLog.Config{
			Format:     "${pid} ${status} - ${method} ${path}\n",
			TimeFormat: time.DateTime,
			TimeZone:   "Asia/Jakarta",
		}))

		fmt.Printf(
			"Allow Origins [%s] & Headers [%s] & Expose Header [%s]",
			config.AllowOrigins,
			config.AllowHeaders,
			config.ExposeHeaders,
		)
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:  config.AllowOrigins,
		AllowHeaders:  config.AllowHeaders,
		ExposeHeaders: config.ExposeHeaders,
	}))

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: config.CookieEncryptKey,
		Except: []string{
			config.CookieKeyName,
		},
	}))

	src.RegisterRoute(app, handler, jwtMiddleware)

	if err := app.Listen(fmt.Sprintf("%s:%s", config.Http.Host, config.Http.Port)); err != nil {
		log.Fatal(err)
	}
}
