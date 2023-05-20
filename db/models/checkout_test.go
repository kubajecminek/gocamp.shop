package models

import (
	"net/url"
	"reflect"
	"testing"

	"gocamp.shop/validator"
)

func TestDecodeCheckout(t *testing.T) {
	b := Billing{}
	p := []Participant{}

	testCases := []struct {
		name string
		in   string
		out  Checkout
	}{
		{"empty", "", Checkout{Billing: b, Participants: p}},
		{
			"allGood",
			"name=John&email=some@email.com&mnumber=111222333",
			Checkout{Name: "John", Email: "some@email.com", MobileNumber: "111222333", Billing: b, Participants: p},
		}}

	for _, tc := range testCases {
		form, _ := url.ParseQuery(tc.in)
		checkout := decodeCheckout(form)
		if !reflect.DeepEqual(checkout, tc.out) {
			t.Errorf("Test - %s: expected %v, actual %v", tc.name, tc.out, checkout)
		}
	}
}

func TestValidate(t *testing.T) {
	v := validator.New()
	testCases := []struct {
		name     string
		input    Checkout
		expected error
	}{
		{
			"CheckoutBug1",
			Checkout{
				Name:         "John Doe",
				MobileNumber: "1234567890",
				Email:        "",
				Billing: Billing{
					Name:    "Johnny Surname",
					Street:  "123 Main St",
					City:    "New York",
					ZipCode: "10001",
					ID:      "",
					TaxID:   "",
				},
				Participants: []Participant{
					Participant{
						Suffix:         "-1-1",
						OrdinalPos:     1,
						ItemID:         1,
						Name:           "John",
						ID:             "1234567890",
						Club:           "",
						WeightCategory: "",
						Belt:           "",
						ShirtSize:      "M",
					},
				},
			},
			validator.ErrZeroValue,
		},
		{
			"CheckoutEmpty",
			Checkout{},
			validator.ErrZeroValue,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.input.Validate(v, 1)
			if err != tc.expected {
				t.Errorf("case: %s, got error %v", tc.name, err)
			}
		})
	}
}

func TestCheckoutInnerEmpty(t *testing.T) {
	in := struct{ c Checkout }{}
	if in.c.Validate(validator.New(), 0) != validator.ErrZeroValue {
		t.Errorf("expected ErrZeroValue")
	}
}
