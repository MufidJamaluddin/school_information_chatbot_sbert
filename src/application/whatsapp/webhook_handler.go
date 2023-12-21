package whatsapp

import (
	"bytes"
	"chatbot_be_go/src/application/whatsapp/dto"
	"chatbot_be_go/src/persistence/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"github.com/gofiber/fiber/v2"
)

type IWhatsAppIntegration interface {
	GetWebhook(c *fiber.Ctx) error
	PostWebhook(c *fiber.Ctx) error
}

type whatsappIntegration struct {
	appConfig  *config.AppConfig
	validator  *validator.Validate
	logger     *logrus.Logger
	httpClient *http.Client
}

// GetWebhook godoc
// @Tags "UC01 Bertanya"
// @Summary Webhook Handle WhatsApp Subscription Ack
// @Description Validate the 'hub.mode' must be 'subscribe' and 'hub.verify_token' must be the same
// @Param hub.mode query string true "Hub Mode" default(subscribe)
// @Param hub.verify_token query string true "Verify Token" default(12321)
// @Param hub.challenge query string true "Challenge" default(TEAM2BINUS)
// @Success 200 {object} string
// @Failure 400 {object} string
// @Router /api/webhook [get]
func (w *whatsappIntegration) GetWebhook(c *fiber.Ctx) error {
	hubMode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if hubMode == "subscribe" && token == w.appConfig.WhatsAppConf.VerifyToken {
		return c.Status(fiber.StatusOK).SendString(challenge)
	}

	return c.SendStatus(fiber.StatusNotFound)
}

// PostWebhook godoc
// @Tags "UC01 Bertanya"
// @Summary Webhook Handle Incoming Message from WhatsApp
// @Description Handle the Incoming Message from WhatsApp API
// @Accept json
// @Produce json
// @Param q body dto.WebhookInDto true "Incoming Message"
// @Success 202 {object} dto.WASendMessageDto
// @Failure 400 {object} string
// @Router /api/webhook [post]
func (w *whatsappIntegration) PostWebhook(c *fiber.Ctx) error {
	var messageDataDto *dto.WebhookInDto
	var messages []*dto.MessageDTO

	mapPhoneNoName := map[string]string{}

	if err := json.Unmarshal(c.Body(), &messageDataDto); err != nil {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Format")
	}

	if messageDataDto.Object != "whatsapp_business_account" {
		return c.Status(fiber.StatusNotAcceptable).SendString("Invalid Object")
	}

	for _, entry := range messageDataDto.Entry {
		for _, change := range entry.Changes {
			for _, contact := range change.Value.Contacts {
				mapPhoneNoName[contact.WaId] = contact.Profile.Name
			}

			for _, message := range change.Value.Messages {
				messages = append(messages, &dto.MessageDTO{
					FromId:       change.Value.Metadata.PhoneNoId,
					FromPhoneNo:  message.From,
					FromFullName: mapPhoneNoName[message.From],
					TextMessage:  message.Text.Body,
					ChatAt:       time.Unix(0, message.Timestamp*int64(time.Millisecond)),
				})
			}
		}
	}

	for _, message := range messages {

		locUrl := fmt.Sprintf(
			"https://graph.facebook.com/v1/messages",
			message.FromId,
		)

		resp := dto.CreateTextWAMessage(
			"",
			message.FromPhoneNo,
			fmt.Sprintf("SERVER: ACK: %s", message.TextMessage),
			false,
		)

		respJson, _ := json.Marshal(&resp)
		reqPath, _ := url.Parse(locUrl)

		buffer := io.NopCloser(bytes.NewBuffer(respJson))

		req := &http.Request{
			Method: "POST",
			URL:    reqPath,
			Header: map[string][]string{
				"Content-Type":  {"application/json"},
				"Authorization": {fmt.Sprintf("Bearer %s", w.appConfig.WhatsAppConf.AccessToken)},
			},
			Body: buffer,
		}

		post, err := w.httpClient.Do(req)
		if err != nil {
			w.logger.Error(err)
			continue
		}

		if post.StatusCode > 299 || post.StatusCode < 199 {
			return c.Status(post.StatusCode).SendStream(post.Body)
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

func New(
	appConfig *config.AppConfig,
	validator *validator.Validate,
	logger *logrus.Logger,
	httpClient *http.Client,
) IWhatsAppIntegration {
	return &whatsappIntegration{
		appConfig:  appConfig,
		validator:  validator,
		logger:     logger,
		httpClient: httpClient,
	}
}
