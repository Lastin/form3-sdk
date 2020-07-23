package form3_sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
)

type SessionCofig struct {
	ApiHost *string
	ApiKey  *string
}

type SdkClient struct {
	config SessionCofig
}

func New(config SessionCofig) *SdkClient {
	if config.ApiHost == nil {
		config.ApiHost = aws.String("http://localhost:8080")
	}
	return &SdkClient{
		config: config,
	}
}

type Data struct {
	Type           *string
	Id             *string
	OrganisationId *string         `json:"organisation_id"`
	Version        int             `json:"omitempty"`
	Relationships  json.RawMessage `json:"omitempty"`
	Attributes     json.RawMessage
}

type Payload struct {
	Data *Data
}

func (client SdkClient) ProcessResponse(resp []byte) (data *Data, err error) {
	response := new(Payload)
	if err = json.Unmarshal(resp, response); err != nil {
		return
	}
	if response.Data == nil {
		return nil, errors.New("data field not present in the payload")
	}
	return response.Data, nil
}

func (client SdkClient) Create(path string, msgType string, attributes json.RawMessage) (data *Data, err error) {
	var reqB, respB []byte
	payload := Payload{
		Data: &Data{
			Type:           &msgType,
			OrganisationId: aws.String(uuid.New().String()),
			Id:             aws.String(uuid.New().String()),
			Attributes:     attributes,
		},
	}
	if reqB, err = json.Marshal(payload); err != nil {
		return
	}
	url := fmt.Sprintf("%s/%s", *client.config.ApiHost, path)
	resp, err := http.Post(url, "application/vnd.api+json", bytes.NewBuffer(reqB))
	if err != nil {
		return
	}
	respB, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return client.ProcessResponse(respB)
}
