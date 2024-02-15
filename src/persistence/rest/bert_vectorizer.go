package rest

import (
	"chatbot_be_go/src/persistence/config"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	resty "github.com/go-resty/resty/v2"
	gerr "github.com/pkg/errors"
)

type ISBertVectorizer interface {
	Encode(text string) ([]float64, error)
}

type sbertVectorizer struct {
	requestor       *resty.Request
	IsFineTuneModel bool
}

func NewSBertVectorizer(config *config.AppConfig) ISBertVectorizer {
	client := resty.New()

	requestor := client.
		SetBaseURL(config.EmbeddingServerHost).
		R().
		SetHeader("Accept", "application/json")

	return &sbertVectorizer{
		requestor:       requestor,
		IsFineTuneModel: config.IsFineTuneModel,
	}
}

func (s *sbertVectorizer) Encode(text string) (result []float64, err error) {
	var response *resty.Response
	var pathStr string

	text = url.QueryEscape(strings.ReplaceAll(text, "/", " "))

	if s.IsFineTuneModel {
		pathStr = fmt.Sprintf("/encode-ft/%s", text)
	} else {
		pathStr = fmt.Sprintf("/encode-ori/%s", text)
	}

	response, err = s.requestor.Get(pathStr)
	if err != nil {
		return
	}

	if response.IsSuccess() {
		err = json.Unmarshal(response.Body(), &result)
	} else {
		err = gerr.Wrap(errors.New(pathStr), response.String())
	}
	return
}
