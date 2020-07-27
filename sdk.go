package form3_sdk

import (
	"bytes"
	"encoding/json"
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
	OrganisationId *string `json:"organisation_id"`
	Version        int
	Relationships  json.RawMessage
	Attributes     json.RawMessage
}

type Payload struct {
	Data *Data
}

func (client SdkClient) Create(path string, msgType string, attributes json.RawMessage) (respB []byte, err error) {
	payload := Payload{
		Data: &Data{
			Type:           &msgType,
			OrganisationId: aws.String(uuid.New().String()),
			Id:             aws.String(uuid.New().String()),
			Attributes:     attributes,
		},
	}
	var reqB []byte
	if reqB, err = json.Marshal(payload); err == nil {
		url := fmt.Sprintf("%s/%s", *client.config.ApiHost, path)
		var resp *http.Response
		if resp, err = http.Post(url, "application/vnd.api+json", bytes.NewBuffer(reqB)); err == nil {
			respB, err = ioutil.ReadAll(resp.Body)
		}
	}
	return
}

func (client SdkClient) Fetch(path, id string) (respB []byte, err error) {
	url := fmt.Sprintf("%s/%s/%s", *client.config.ApiHost, path, id)
	var resp *http.Response
	if resp, err = http.Get(url); err == nil {
		respB, err = ioutil.ReadAll(resp.Body)
	}
	return
}
