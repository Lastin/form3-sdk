package form3_sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"reflect"
)

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
	return client.RawGet(fmt.Sprintf("%s/%s", path, id))
}

func buildFilter(i interface{}) string {
	if i == nil {
		return ""
	}
	var buf bytes.Buffer
	v := reflect.ValueOf(i)
	for j := 0; j < v.NumField(); j++ {
		if v.Field(j).IsValid() && !v.Field(j).IsZero() {
			jsonTagName := v.Type().Field(j).Tag.Get("json")
			if len(jsonTagName) > 0 {
				buf.WriteString(fmt.Sprintf("&filter[%s]=%v", jsonTagName, reflect.Indirect(v.Field(j)).Interface()))
			}

		}
	}
	return buf.String()
}

func (client SdkClient) List(path string, pageNumber, pageSize int, filter interface{}) (respB []byte, err error) {
	//v1/organisation/accounts?page[number]={page_number}&page[size]={page_size}&filter[{attribute}]={filter_value}
	filters := buildFilter(filter)
	return client.RawGet(fmt.Sprintf("%s?page[number]=%d&page[size]=%d%s", path, pageNumber, pageSize, filters))
}

func (client SdkClient) RawGet(path string) (respB []byte, err error) {
	var req *http.Request
	url := fmt.Sprintf("%s/%s", *client.config.ApiHost, path)
	if req, err = http.NewRequest("GET", url, nil); err == nil {
		return client.Do(req)
	}
	return
}

func (client SdkClient) Delete(path, id string, version int) (respB []byte, err error) {
	url := fmt.Sprintf("%s/%s/%s?version=%d", *client.config.ApiHost, path, id, version)
	var req *http.Request
	if req, err = http.NewRequest("DELETE", url, nil); err == nil {
		return client.Do(req)
	}
	return
}

func (client SdkClient) Do(req *http.Request) (respB []byte, err error) {
	var resp *http.Response
	if resp, err = client.httpClient.Do(req); err == nil {
		respB, err = ioutil.ReadAll(resp.Body)
	}
	return
}
