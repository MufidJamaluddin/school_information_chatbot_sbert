package chat_user

import (
	"chatbot_be_go/src/application/shared"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type IChatUserHandler interface {
	ListChatUser(c *fiber.Ctx) error
}

type chatUserHandler struct {
	logger             *logrus.Logger
	chatUserRepository IChatUserRepository
}

func NewChatUserHandler(
	logger *logrus.Logger,
	chatUserRepository IChatUserRepository,
) IChatUserHandler {
	return &chatUserHandler{
		logger:             logger,
		chatUserRepository: chatUserRepository,
	}
}

var _ IChatUserHandler = &chatUserHandler{}

// ListChatUser godoc
// @Security BasicAuth
// @securityDefinitions.basic BasicAuth
// @Tags "UC05 Menjawab Pertanyaan Manual"
// @Summary Get the List of Chat User Data
// @Description Show List of Chat User Data
// @Param keyword query string false "keyword"
// @Param start query string true "start"
// @Param size query string true "size"
// @Success 200 {object} []dto.ChatUserItemDTO
// @Failure 400 {object} string
// @Router /api/chat-users [get]
func (cu *chatUserHandler) ListChatUser(c *fiber.Ctx) error {
	var counter uint64

	keyword := c.Query("keyword")

	start, _ := strconv.ParseUint(c.Query("start"), 10, 32)
	size, _ := strconv.ParseUint(c.Query("size"), 10, 32)

	if size == 0 {
		size = 10
	}

	callbackSendData := shared.CreateSendCallback(start, &counter, c)

	_, _ = c.Write([]byte("["))

	if totalAll, err := cu.chatUserRepository.ListChatUser(
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
