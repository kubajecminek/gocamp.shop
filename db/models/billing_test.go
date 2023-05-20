package models

import (
	"net/url"
	"testing"

	"gocamp.shop/validator"
)

func TestBillingValid(t *testing.T) {
	v := validator.New()
	testCases := []struct {
		name string
		in   Billing
		out  error
	}{
		{"empty", Billing{}, validator.ErrZeroValue},
		{"missingOne", Billing{Name: "John", Street: "Downing Street 2", City: "London", ZipCode: "", ID: "", TaxID: ""}, validator.ErrZeroValue},
		//		{"wrongID", Billing{Name: "John", Street: "Downing Street 2", City: "London", ZipCode: "1234 56", ID: "ABC", TaxID: "CZ1234"}, validator.ErrRegexp},
		{"allGood", Billing{Name: "John", Street: "Downing Street 2", City: "London", ZipCode: "1234 56", ID: "123", TaxID: "CZ1234"}, nil},
		{"maxLenExceeded", Billing{Name: "1234567890123456789012345678901"}, validator.ErrMaxLen},
	}

	for _, tc := range testCases {
		err := v.Validate(tc.in)
		if err != tc.out {
			t.Errorf("Test - %s: expected %v, actual %v", tc.name, tc.out, err)
		}
	}
}

func TestDecodeBilling(t *testing.T) {
	testCases := []struct {
		name string
		in   string
		out  Billing
	}{
		{"empty", "", Billing{}},
		{
			"allGood", "bname=John&bstreet=Downing+Street+2&bcity=London&bzipcode=1234+56&binum=123&btnum=CZ1234",
			Billing{Name: "John", Street: "Downing Street 2", City: "London", ZipCode: "1234 56", ID: "123", TaxID: "CZ1234"}},
		{
			"missingOne", "bname=John&bstreet=Downing+Street+2&bcity=London&bzipcode=1234+56&binum=&btnum=CZ1234",
			Billing{Name: "John", Street: "Downing Street 2", City: "London", ZipCode: "1234 56", ID: "", TaxID: "CZ1234"}},
	}

	for _, tc := range testCases {
		form, _ := url.ParseQuery(tc.in)
		billing := decodeBilling(form)
		if billing != tc.out {
			t.Errorf("Test - %s: expected %v, actual %v", tc.name, tc.out, billing)
		}
	}
}
