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

const ContentType = "application/vnd.api+json"

type SessionCofig struct {
	ApiHost *string
	ApiKey  *string
}

type SdkClient struct {
	httpClient http.Client
	config     SessionCofig
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

func (client SdkClient) GetPath(path string) (respB []byte, err error) {
	var resp *http.Response
	if resp, err = client.httpClient.Get(fmt.Sprintf("%s/%s", *client.config.ApiHost, path)); err == nil {
		respB, err = ioutil.ReadAll(resp.Body)
	}
	return
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
		if resp, err = client.httpClient.Post(url, ContentType, bytes.NewBuffer(reqB)); err == nil {
			respB, err = ioutil.ReadAll(resp.Body)
		}
	}
	return
}

func (client SdkClient) Fetch(path, id string) (respB []byte, err error) {
	return client.GetPath(fmt.Sprintf("%s/%s", path, id))
}

func (client SdkClient) List(path string, pageNumber, pageSize int, filter interface{}) (respB []byte, err error) {
	filterStr := buildFilter(filter)
	return client.GetPath(fmt.Sprintf(
		"%s?page[number]=%d&page[size]=%d%s",
		path,
		pageNumber,
		pageSize,
		filterStr,
	))
}

func (client SdkClient) Delete(path, id string, version int) (success bool, err error) {
	url := fmt.Sprintf("%s/%s/%s?version=%d", *client.config.ApiHost, path, id, version)
	var req *http.Request
	if req, err = http.NewRequest("DELETE", url, nil); err == nil {
		var resp *http.Response
		if resp, err = client.httpClient.Do(req); err == nil {
			if resp.StatusCode == http.StatusNoContent {
				success = true
			} else if resp.StatusCode == http.StatusNotFound {
				err = errors.New("specified resource does not exist")
			} else if resp.StatusCode == http.StatusConflict {
				err = errors.New("specified version incorrect")
			}
		}
	}
	return
}
