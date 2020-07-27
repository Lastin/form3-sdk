package accounts

import (
	"encoding/json"
	form3_sdk "github.com/Lastin/form3-sdk"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

var a1 = &Account{
	Country:                 aws.String("GB"),
	BaseCurrency:            aws.String("GBP"),
	AccountNumber:           aws.String("41426819"),
	BankId:                  aws.String("400300"),
	BankIdCode:              aws.String("GBDSC"),
	Bic:                     aws.String("NWBKGB22"),
	IBan:                    aws.String("GB11NWBK40030041426819"),
	Name:                    []*string{aws.String("Samantha Holder")},
	AlternativeNames:        []*string{aws.String("Sam Holder")},
	AccountClassification:   aws.String("Personal"),
	JointAccount:            aws.Bool(false),
	AccountMatchingOptOut:   aws.Bool(false),
	SecondaryIdentification: aws.String("A1B2C3D4"),
	Switched:                aws.Bool(false),
	PrivateIdentification: PrivateIdentification{
		BirthDate:      aws.String("2017-07-23"),
		BirthCountry:   aws.String("GB"),
		Identification: aws.String("13YH458762"),
		Address:        aws.String("[10 Avenue des Champs]"),
		City:           aws.String("London"),
		Country:        aws.String("GB"),
	},
	OrganisationIdentification: OrganisationIdentification{
		Identification: aws.String("123654"),
		Actors: []*Actor{&Actor{
			Name:      []*string{aws.String("Jeff Page")},
			BirthDate: aws.String("1970-01-01"),
			Residency: aws.String("GB"),
		}},
		Address: []*string{aws.String("10 Avenue des Champs")},
		City:    aws.String("London"),
		Country: aws.String("GB"),
	},
}

//note the discrepancy between address in Private and Organisational Identification
func Test_JSONUnmarshal(t *testing.T) {
	b1, err := ioutil.ReadFile("../test_files/account_response.json")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name            string
		responseBytes   []byte
		expectedAccount *Account
		wantErr         bool
	}{
		{"valid response", b1, a1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := new(Create)
			err := json.Unmarshal(tt.responseBytes, result)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tt.expectedAccount, result.Data.Attributes)
			}
		})
	}
}

func Test_Create(t *testing.T) {
	result, err := New(form3_sdk.SessionCofig{}).Create(a1)
	assert.NoError(t, err)
	assert.EqualValues(t, &Account{
		Country:                 a1.Country,
		BaseCurrency:            a1.BaseCurrency,
		AccountNumber:           a1.AccountNumber,
		BankId:                  a1.BankId,
		BankIdCode:              a1.BankIdCode,
		Bic:                     a1.Bic,
		IBan:                    a1.IBan,
		AccountClassification:   a1.AccountClassification,
		JointAccount:            a1.JointAccount,
		AccountMatchingOptOut:   a1.AccountMatchingOptOut,
		SecondaryIdentification: a1.SecondaryIdentification,
	}, result.Data.Attributes)
}

func TestAccounts_Fetch(t *testing.T) {
	result, err := New(form3_sdk.SessionCofig{}).Create(a1)
	assert.NoError(t, err)
	result, err = New(form3_sdk.SessionCofig{}).Fetch(*result.Data.Id)
	assert.NoError(t, err)
	assert.EqualValues(t, &Account{
		Country:                 a1.Country,
		BaseCurrency:            a1.BaseCurrency,
		AccountNumber:           a1.AccountNumber,
		BankId:                  a1.BankId,
		BankIdCode:              a1.BankIdCode,
		Bic:                     a1.Bic,
		IBan:                    a1.IBan,
		AccountClassification:   a1.AccountClassification,
		JointAccount:            a1.JointAccount,
		AccountMatchingOptOut:   a1.AccountMatchingOptOut,
		SecondaryIdentification: a1.SecondaryIdentification,
	}, result.Data.Attributes)
}

func TestAccounts_List(t *testing.T) {
	client := New(form3_sdk.SessionCofig{})
	tests := []struct {
		name       string
		pageNumber int
		pageSize   int
	}{
		{"one", 0, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.List(tt.pageNumber, tt.pageSize, Account{
				AccountNumber: aws.String("41426819"),
			})
			assert.NoError(t, err)
			//list.Next()
		})
	}
}
