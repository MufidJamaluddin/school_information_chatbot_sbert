package rest

import (
	"chatbot_be_go/src/persistence/config"
	"encoding/json"
	"fmt"

	resty "github.com/go-resty/resty/v2"
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
		SetTimeout(config.Http.S2STimeout).
		SetHostURL(config.EmbeddingServerHost).
		R().
		SetHeader("Accept", "application/json")

	return &sbertVectorizer{
		requestor:       requestor,
		IsFineTuneModel: config.IsFineTuneModel,
	}
}

func (s *sbertVectorizer) Encode(text string) (result []float64, err error) {
	var response *resty.Response

	if s.IsFineTuneModel {
		response, err = s.requestor.Get(fmt.Sprintf("/encode-ft/%s", text))
	} else {
		response, err = s.requestor.Get(fmt.Sprintf("/encode-ori/%s", text))
	}

	if err != nil {
		return
	}

	err = json.Unmarshal(response.Body(), &result)
	return
}
