package accounts

import (
	"encoding/json"
	form3_sdk "github.com/Lastin/form3-sdk"
)

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

func processResponse(resp []byte) (*Account, error) {
	account := new(Account)
	if err := json.Unmarshal(resp, account); err != nil {
		return nil, err
	}
	return account, nil
}

const path = "v1/organisation/accounts"
const MsgType = "accounts"

func (client Accounts) Create(account *Account) (*Account, error) {
	if accountB, err := json.Marshal(account); err != nil {
		return nil, err
	} else {
		resp, err := client.sdkClient.Create(path, MsgType, accountB)
		if err != nil {
			return nil, err
		}
		return processResponse(resp.Attributes)
	}
}

func Fetch()  {}
func List()   {}
func Delete() {}
