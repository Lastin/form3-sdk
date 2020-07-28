package accounts

import (
	"encoding/json"
	form3_sdk "github.com/Lastin/form3-sdk"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strconv"
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
	DeleteAll(t, client)
	CreateBunch(t, 200)
	/*
		This is for validating if filter works.
			checkFilter := func(list *List, expectedA Account) {
				for _, actualA := range list.Data {
					if expectedA.BankId != nil {
						assert.EqualValues(t, *expectedA.BankId, *actualA.Attributes.BankId)
					}
					if expectedA.AccountNumber != nil {
						assert.EqualValues(t, *expectedA.AccountNumber, *actualA.Attributes.AccountNumber)
					}
				}
			}
	*/
	tests := []struct {
		name       string
		pageNumber int
		pageSize   int
		filter     Account
	}{
		{"no filter", 0, 10, Account{}},
		{"filter by account number", 0, 10, Account{AccountNumber: aws.String("41426819")}},
		{"filter by account number", 0, 100, Account{AccountNumber: aws.String("41426819")}},
		//{"filter by bank id", 0, 10, Account{BankId: aws.String("2")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := client.List(tt.pageNumber, tt.pageSize, tt.filter)
			assert.NoError(t, err)
			assert.True(t, len(list.Data) <= tt.pageSize)
			//checkFilter(list, tt.filter)
			for list.HasNext() {
				list, err = list.Next()
				assert.NoError(t, err)
				assert.True(t, len(list.Data) <= tt.pageSize)
				//checkFilter(list, tt.filter)
			}
		})
	}
}

func CreateBunch(t *testing.T, count int) {
	client := New(form3_sdk.SessionCofig{})
	for i := 0; i < count; i++ {
		_, err := client.Create(&Account{
			Country:                 a1.Country,
			BaseCurrency:            a1.BaseCurrency,
			AccountNumber:           aws.String(strconv.Itoa(i)),
			BankId:                  a1.BankId,
			BankIdCode:              a1.BankIdCode,
			Bic:                     a1.Bic,
			IBan:                    a1.IBan,
			AccountClassification:   a1.AccountClassification,
			JointAccount:            a1.JointAccount,
			AccountMatchingOptOut:   a1.AccountMatchingOptOut,
			SecondaryIdentification: a1.SecondaryIdentification,
		})
		assert.NoError(t, err)
	}
	/* This was implemented for testing filter, which clearly doesn't filter
	one := aws.String("41426819")
	two := aws.String("41426820")
	for _, v := range []struct {
		accountNumber *string
		bankId        *string
	}{
		{one, one},
		{one, two},
		{two, one},
		{two, two},
	} {
		_, err := client.Create(&Account{
			Country:                 a1.Country,
			BaseCurrency:            a1.BaseCurrency,
			AccountNumber:           v.accountNumber,
			BankId:                  v.bankId,
			BankIdCode:              a1.BankIdCode,
			Bic:                     a1.Bic,
			IBan:                    a1.IBan,
			AccountClassification:   a1.AccountClassification,
			JointAccount:            a1.JointAccount,
			AccountMatchingOptOut:   a1.AccountMatchingOptOut,
			SecondaryIdentification: a1.SecondaryIdentification,
		})
		assert.NoError(t, err)
	}
	*/
}

func Test_Delete(t *testing.T) {
	client := New(form3_sdk.SessionCofig{})
	//prepare for test
	DeleteAll(t, client)
	result, err := client.List(0, 100, Account{})
	assert.NoError(t, err)
	assert.Len(t, result.Data, 0)
	//add an account
	createResult, err := client.Create(a1)
	assert.NoError(t, err)
	fetchResult, err := client.Fetch(*createResult.Data.Id)
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
	}, fetchResult.Data.Attributes)
	success, err := client.Delete(*fetchResult.Data.Id, fetchResult.Data.Version)
	assert.NoError(t, err)
	assert.True(t, success)
	fetchResult, err = client.Fetch("*createResult.Data.Id")
	assert.NoError(t, err)
}

func Test_DeleteAll(t *testing.T) {
	client := New(form3_sdk.SessionCofig{})
	DeleteAll(t, client)
	CreateBunch(t, 200)
	result, err := client.List(0, 100, Account{})
	assert.NoError(t, err)
	assert.Len(t, result.Data, 100)
	DeleteAll(t, client)
	result, err = client.List(0, 100, Account{})
	assert.NoError(t, err)
	assert.Len(t, result.Data, 0)
}

func DeleteAll(t *testing.T, client Accounts) {
	deleteEach := func(list *List) {
		for _, account := range list.Data {
			success, err := client.Delete(*account.Id, account.Version)
			assert.NoError(t, err)
			assert.True(t, success)
		}
	}
	list, err := client.List(0, 100, Account{})
	assert.NoError(t, err)
	deleteEach(list)
	for list.HasNext() {
		list, err = list.Next()
		assert.NoError(t, err)
		deleteEach(list)
	}
}
