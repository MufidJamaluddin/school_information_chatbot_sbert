package rest

import (
	"bytes"
	"chatbot_be_go/src/persistence/config"
	"chatbot_be_go/src/persistence/rest/dto"
	"chatbot_be_go/src/persistence/rest/factory"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
)

type IWAMessageAPI interface {
	SendTextMessage(
		fromWaId string,
		contextMessageId string,
		recipientTo string,
		text string,
		includeUrl bool,
	) error

	SendImageMessage(
		fromWaId string,
		contextMessageId string,
		recipientTo string,
		imageUrl string,
	) error

	SendLocation(
		fromWaId string,
		contextMessageId string,
		recipientTo string,
		locationName string,
		locationAddress string,
		latitude float32,
		longitude float32,
	) error
}

type waMessageAPI struct {
	logger       *logrus.Logger
	httpClient   *http.Client
	whatsappConf config.WhatsAppConf
}

func (w *waMessageAPI) sendMessage(fromId string, message *dto.WASendMessageDto) error {
	locUrl := fmt.Sprintf(
		w.whatsappConf.MessageAPIPattern,
		fromId,
	)

	respJson, err := json.Marshal(&message)
	if err != nil {
		return err
	}

	reqPath, err := url.Parse(locUrl)
	if err != nil {
		return err
	}

	buffer := io.NopCloser(bytes.NewBuffer(respJson))

	req := &http.Request{
		Method: "POST",
		URL:    reqPath,
		Header: map[string][]string{
			"Content-Type":  {"application/json"},
			"Authorization": {fmt.Sprintf("Bearer %s", w.whatsappConf.AccessToken)},
		},
		Body: buffer,
	}

	_, err = w.httpClient.Do(req)

	if err != nil {
		w.logger.Error("Error in Send Message to WhatsApp API", err)
	}

	return err
}

func (w *waMessageAPI) SendTextMessage(
	fromWaId string,
	contextMessageId string,
	recipientTo string,
	text string,
	includeUrl bool,
) error {
	textMsg := factory.CreateTextWAMessage(
		contextMessageId,
		recipientTo,
		text,
		includeUrl,
	)

	return w.sendMessage(fromWaId, textMsg)
}

func (w *waMessageAPI) SendImageMessage(
	fromWaId string,
	contextMessageId string,
	recipientTo string,
	imageUrl string,
) error {
	imageMsg := factory.CreateImageWAMessage(
		contextMessageId,
		recipientTo,
		imageUrl,
	)

	return w.sendMessage(fromWaId, imageMsg)
}

func (w *waMessageAPI) SendLocation(
	fromWaId string,
	contextMessageId string,
	recipientTo string,
	locationName string,
	locationAddress string,
	latitude float32,
	longitude float32,
) error {
	locationMsg := factory.CreateLocationWAMessage(
		contextMessageId,
		recipientTo,
		locationName,
		locationAddress,
		latitude,
		longitude,
	)

	return w.sendMessage(fromWaId, locationMsg)
}

func NewWAMessageAPI(
	logger *logrus.Logger,
	httpClient *http.Client,
	whatsappConf config.WhatsAppConf,
) IWAMessageAPI {
	return &waMessageAPI{
		logger:       logger,
		httpClient:   httpClient,
		whatsappConf: whatsappConf,
	}
}
