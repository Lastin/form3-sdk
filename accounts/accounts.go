package accounts

import (
	"encoding/json"
	"github.com/Lastin/form3-sdk"
)

const (
	apiPath = "v1/organisation/accounts"
	MsgType = "accounts"
)

// Root struct of this package. Provides range of functions to interact with accounts API
type Accounts struct {
	sdkClient *form3.SdkClient
}

// Represents type of function that needs to be provided to the iterator
type IteratorFunc func(i int, data *AccountData) error

// Represents the response received from the API
type AccountsResponse struct {
	Data  *AccountData
	Links Links
}

// Represents the data segment of the AccountsResponse. It is expected to be generic across other API endpoints
type AccountData struct {
	Attributes     *Account `json:"attributes"`
	CreatedOn      *string  `json:"created_on"`
	Id             *string
	ModifiedOn     *string `json:"modified_on"`
	OrganisationId *string `json:"organisation_id"`
	Type           *string
	Version        int
}

// Represents structure of Links element of AccountsResponse. Provides easy access to pagination function
type Links struct {
	First *string
	Last  *string
	Next  *string
	Self  *string
}

// Represents response received by the API when fetching list of accounts
// Also carries pointer to instance of *form3_sdk.SdkClient used by Next() and First() functions
type List struct {
	sdkClient *form3.SdkClient
	Data      []*AccountData
	Links     *Links
}

// Represents actual Account data provided in the "Attributes" field of the AccountsResponse
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

// Represents the structure element "PrivateIdentification" of Account structure
type PrivateIdentification struct {
	BirthDate      *string `json:"birth_date"`
	BirthCountry   *string `json:"birth_country"`
	Identification *string
	Address        *string
	City           *string
	Country        *string
}

// Represents the structure element "OrganisationIdentification" of Account structure
type OrganisationIdentification struct {
	Identification *string
	Actors         []*Actor
	Address        []*string
	City           *string
	Country        *string
}

// Represents the structure element "Actor" of OrganisationIdentification structure
type Actor struct {
	Name      []*string
	BirthDate *string `json:"birth_date"`
	Residency *string
}

// Creates a new instance of Accounts which carries instance of form3_sdk.SdkClient with config as provided in the argument
func New(config form3.SessionCofig) *Accounts {
	sdkClient := form3.New(config)
	return &Accounts{sdkClient: sdkClient}
}

// Sends request to create an account reflecting provided account object
// On success the returned data is populated into AccountsResponse
func (client Accounts) Create(account *Account) (result *AccountsResponse, err error) {
	var reqB []byte
	if reqB, err = json.Marshal(account); err == nil {
		var respB []byte
		if respB, err = client.sdkClient.RestMakePostRequest(apiPath, MsgType, reqB); err == nil {
			result = new(AccountsResponse)
			err = json.Unmarshal(respB, result)
		}
	}
	return
}

// Fetches single account data with given account id
func (client Accounts) Fetch(id string) (result *AccountsResponse, err error) {
	data, err := client.sdkClient.Fetch(apiPath, id)
	if err == nil {
		result = new(AccountsResponse)
		err = json.Unmarshal(data, result)
	}
	return
}

// Fetches and returns list of accounts, where page number and it's size are as provided
// Also can filter by given account attributes
func (client Accounts) List(pageNumber, pageSize int, filter Account) (result *List, err error) {
	var data []byte
	if data, err = client.sdkClient.List(apiPath, pageNumber, pageSize, filter); err == nil {
		result = new(List)
		result.sdkClient = client.sdkClient
		err = json.Unmarshal(data, result)
	}
	return
}

// Fetches and returns first page of the list
func (list List) First() (first *List, err error) {
	first = new(List)
	first.sdkClient = list.sdkClient
	var data []byte
	if data, err = list.sdkClient.RestMakeGetRequest(*list.Links.First); err == nil {
		err = json.Unmarshal(data, first)
	}
	return
}

// Fetches and returns the next page of the list
func (list List) Next() (next *List, err error) {
	next = new(List)
	next.sdkClient = list.sdkClient
	if list.HasNext() {
		var data []byte
		if data, err = list.sdkClient.RestMakeGetRequest(*list.Links.Next); err == nil {
			err = json.Unmarshal(data, next)
		}
	}
	return
}

// Returns true if the list has link to next page
func (list *List) HasNext() bool {
	return list.Links != nil && list.Links.Next != nil
}

func (list *List) Iterate(f IteratorFunc) (err error) {
	for i, account := range list.Data {
		if err = f(i, account); err != nil {
			return
		}
	}
	return
}

// Given a list iterates over each element until no more elements available
func (list *List) Walk(f IteratorFunc) (err error) {
	if err = list.Iterate(f); err != nil {
		return
	}
	for list.HasNext() {
		if list, err = list.Next(); err != nil {
			return
		}
		if err = list.Iterate(f); err != nil {
			return
		}
	}
	return
}

// Deletes an account with given id and version
func (client Accounts) Delete(id string, version int) (bool, error) {
	return client.sdkClient.Delete(apiPath, id, version)
}
