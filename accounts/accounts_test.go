package accounts

import (
	"encoding/json"
	"fmt"
	"github.com/Lastin/form3-sdk"
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
		Actors: []*Actor{
			{
				Name:      []*string{aws.String("Jeff Page")},
				BirthDate: aws.String("1970-01-01"),
				Residency: aws.String("GB"),
			},
		},
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
			result := new(AccountsResponse)
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

// Tests creating of the account
func Test_Create(t *testing.T) {
	result, err := New(form3.SessionCofig{}).Create(a1)
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

// Tests fetching of a freshly created account
func TestAccounts_Fetch(t *testing.T) {
	result, err := New(form3.SessionCofig{}).Create(a1)
	assert.NoError(t, err)
	result, err = New(form3.SessionCofig{}).Fetch(*result.Data.Id)
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

// Tests getting a list of accounts
func TestAccounts_List(t *testing.T) {
	client := New(form3.SessionCofig{})
	DeleteAll(t, client)
	CreateBunch(t, 250)
	/* This is for validating if filter works.
	checkFilter := func (list *List, expectedA Account) {
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
		name            string
		startPageNumber int
		pageSize        int
		filter          Account
	}{
		//{"filter by account number", 0, 10, Account{AccountNumber: aws.String("41426819")}},
		{"no filter", 0, 10, Account{}},
		{"page size 10", 0, 10, Account{}},
		{"page size 20", 0, 20, Account{}},
		{"starting page 2 and page size 20", 1, 20, Account{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pageNumber := tt.startPageNumber
			list, err := client.List(pageNumber, tt.pageSize, tt.filter)
			assert.NoError(t, err)
			assert.LessOrEqual(t, len(list.Data), tt.pageSize)
			for i, each := range list.Data {
				expectedAccountNumber := strconv.Itoa(i + (tt.pageSize * pageNumber))
				assert.Equal(t, expectedAccountNumber, *each.Attributes.AccountNumber)
			}
			//checkFilter(list, tt.filter)
			for list.HasNext() {
				pageNumber++
				list, err = client.List(pageNumber, tt.pageSize, tt.filter)
				assert.NoError(t, err)
				assert.LessOrEqual(t, len(list.Data), tt.pageSize)
				for i, each := range list.Data {
					expectedAccountNumber := strconv.Itoa(i + (tt.pageSize * pageNumber))
					assert.Equal(t, expectedAccountNumber, *each.Attributes.AccountNumber)
				}
				//checkFilter(list, tt.filter)
			}
		})
	}
}

// Test walking of the list until the end
func TestList_Walk(t *testing.T) {
	client := New(form3.SessionCofig{})
	DeleteAll(t, client)
	CreateBunch(t, 250)
	list, err := client.List(0, 100, Account{})
	assert.NoError(t, err)
	total := 0
	err = list.Walk(func(_ int, accountData *AccountData) error {
		total++
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, 250, total)
}

func Test_Delete(t *testing.T) {
	client := New(form3.SessionCofig{})
	DeleteAll(t, client)
	result, err := client.List(0, 100, Account{})
	assert.NoError(t, err)
	assert.Len(t, result.Data, 0)
	// add an account
	createResult, err := client.Create(a1)
	assert.NoError(t, err)
	fetchResult, err := client.Fetch(*createResult.Data.Id)
	assert.NoError(t, err)
	// assert account exists
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
	// now delete and check
	success, err := client.Delete(*fetchResult.Data.Id, fetchResult.Data.Version)
	assert.NoError(t, err)
	assert.True(t, success)
	fetchResult, err = client.Fetch("*createResult.Data.Id")
	assert.NoError(t, err)
}

// Helper func for creating desired number of dummy accounts with incremental Ids
func CreateBunch(t *testing.T, count int) {
	client := New(form3.SessionCofig{})
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
		_, err := client.RestMakePostRequest(&Account{
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

/* Helper function for clearing the accounts
Actually pagination is not working correctly at backend, as proven in given scenario:
- add 250 accounts
- query first 100
- delete first 100
- use provided link to query next 100
- we delete next 100
- use provided link to query next 100 (50 left in this case), and we get empty list because our accounts shifted left

Therefore here we keep querying first page all the time
*/
func DeleteAll(t *testing.T, client *Accounts) {
	count := 0
	deleteFunc := func(_ int, accountData *AccountData) error {
		success, err := client.Delete(*accountData.Id, accountData.Version)
		assert.True(t, success)
		count++
		return err
	}
	for list, err := client.List(0, 100, Account{}); err == nil && len(list.Data) > 0; list, err = list.First() {
		err = list.Iterate(deleteFunc)
		assert.NoError(t, err)
	}
	fmt.Println("deleted", count)
}
