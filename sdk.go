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

// Sdk client allows for easy interactions with the API
type SdkClient struct {
	HttpClient http.Client
	config     SessionCofig
}

// Represents the structure sent and received from various API endpoints
type Payload struct {
	Data *Data
}

// Represents the "Data" structure sent and returned as part of the Payload
type Data struct {
	Type           *string
	Id             *string
	OrganisationId *string `json:"organisation_id"`
	Version        int
	Relationships  json.RawMessage
	Attributes     json.RawMessage
}

// Creates a new instance of SdkClient and returns it's pointer
func New(config SessionCofig) *SdkClient {
	if config.ApiHost == nil {
		config.ApiHost = aws.String("http://localhost:8080")
	}
	return &SdkClient{
		config: config,
	}
}

// Queries fqdn provided by the config in the client followed by the "path" using GET method and.
func (client SdkClient) RestMakeGetRequest(path string) (respB []byte, err error) {
	var resp *http.Response
	if resp, err = client.HttpClient.Get(fmt.Sprintf("%s/%s", *client.config.ApiHost, path)); err == nil {
		respB, err = ioutil.ReadAll(resp.Body)
	}
	return
}

// Send a POST request to API to (most likely) create a given resource
func (client SdkClient) RestMakePostRequest(path string, msgType string, attributes json.RawMessage) (respB []byte, err error) {
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
		if resp, err = client.HttpClient.Post(url, ContentType, bytes.NewBuffer(reqB)); err == nil {
			respB, err = ioutil.ReadAll(resp.Body)
		}
	}
	return
}

// Fetches a resource with given id from provided endpoint/path
func (client SdkClient) Fetch(path, id string) (respB []byte, err error) {
	return client.RestMakeGetRequest(fmt.Sprintf("%s/%s", path, id))
}

// Fetches a list of resources from provided endpoint/path
func (client SdkClient) List(path string, pageNumber, pageSize int, filter interface{}) (respB []byte, err error) {
	filterStr := buildFilter(filter)
	return client.RestMakeGetRequest(fmt.Sprintf(
		"%s?page[number]=%d&page[size]=%d%s",
		path,
		pageNumber,
		pageSize,
		filterStr,
	))
}

// Deletes a resource with given id and version from provided endpoint/path
func (client SdkClient) Delete(path, id string, version int) (success bool, err error) {
	url := fmt.Sprintf("%s/%s/%s?version=%d", *client.config.ApiHost, path, id, version)
	var req *http.Request
	if req, err = http.NewRequest("DELETE", url, nil); err == nil {
		var resp *http.Response
		if resp, err = client.HttpClient.Do(req); err == nil {
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
