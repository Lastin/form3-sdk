package accounts

import (
	"encoding/json"
	form3_sdk "github.com/Lastin/form3-sdk"
)

type AccountData struct {
	Attributes     *Account `json:"attributes"`
	CreatedOn      *string  `json:"created_on"`
	Id             *string
	ModifiedOn     *string `json:"modified_on"`
	OrganisationId *string `json:"organisation_id"`
	Type           *string
	Version        int
}

type Create struct {
	Data  *AccountData
	Links Links
}

type Links struct {
	First *string
	Last  *string
	Next  *string
	Self  *string
}

type List struct {
	sdkClient *form3_sdk.SdkClient
	Data      []*AccountData
	Links     *Links
}

type Account struct {
	Country                    *string
	BaseCurrency               *string `json:"base_currency"`
	AccountNumber              *string `json:"account_number"`
	BankId                     *string `json:"bank_id"`
	BankIdCode                 *string `json:"bank_id_code"`
	Bic                        *string `json:"bic"`
	IBan                       *string
	Name                       []*string
	AlternativeNames           []*string `json:"alternative_names"`
	AccountClassification      *string   `json:"account_classification"`
	JointAccount               *bool     `json:"joint_account"`
	AccountMatchingOptOut      *bool     `json:"account_matching_opt_out"`
	SecondaryIdentification    *string   `json:"secondary_identification"`
	Switched                   *bool
	PrivateIdentification      PrivateIdentification      `json:"private_identification"`
	OrganisationIdentification OrganisationIdentification `json:"organisation_identification"`
}

type PrivateIdentification struct {
	BirthDate      *string `json:"birth_date"`
	BirthCountry   *string `json:"birth_country"`
	Identification *string
	Address        *string
	City           *string
	Country        *string
}

type OrganisationIdentification struct {
	Identification *string
	Actors         []*Actor
	Address        []*string
	City           *string
	Country        *string
}

type Actor struct {
	Name      []*string
	BirthDate *string `json:"birth_date"`
	Residency *string
}

type Accounts struct {
	sdkClient *form3_sdk.SdkClient
}

func New(config form3_sdk.SessionCofig) Accounts {
	sdkClient := form3_sdk.New(config)
	return Accounts{sdkClient: sdkClient}
}

const (
	createPath = "v1/organisation/accounts"
	fetchPath  = "v1/organisation/accounts"
	listPath   = "v1/organisation/accounts"
	deletePath = "v1/organisation/accounts"
)
const MsgType = "accounts"

func (client Accounts) Create(account *Account) (result *Create, err error) {
	var reqB []byte
	if reqB, err = json.Marshal(account); err == nil {
		var respB []byte
		if respB, err = client.sdkClient.Create(createPath, MsgType, reqB); err == nil {
			result = new(Create)
			err = json.Unmarshal(respB, result)
		}
	}
	return
}

func (client Accounts) Fetch(id string) (result *Create, err error) {
	data, err := client.sdkClient.Fetch(fetchPath, id)
	if err == nil {
		result = new(Create)
		err = json.Unmarshal(data, result)
	}
	return
}

func (client Accounts) List(pageNumber, pageSize int, filter Account) (result *List, err error) {
	var data []byte
	if data, err = client.sdkClient.List(listPath, pageNumber, pageSize, filter); err == nil {
		result = new(List)
		result.sdkClient = client.sdkClient
		err = json.Unmarshal(data, result)
	}
	return
}

func (list List) Next() (result *List, err error) {
	if list.HasNext() {
		var data []byte
		if data, err = list.sdkClient.GetPath(*list.Links.Next); err == nil {
			result = new(List)
			result.sdkClient = list.sdkClient
			err = json.Unmarshal(data, result)
		}
	}
	return
}

func (list List) HasNext() bool {
	return list.Links != nil && list.Links.Next != nil
}

func (client Accounts) Delete(id string, version int) (bool, error) {
	return client.sdkClient.Delete(deletePath, id, version)
}
