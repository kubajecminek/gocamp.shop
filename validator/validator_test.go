package validator

import (
	"reflect"
	"testing"
)

func cmp(a, b []Tag) bool {
	for i, x := range a {
		x.Fn = nil
		b[i].Fn = nil
		if !reflect.DeepEqual(x, b[i]) {
			return false
		}
	}
	return true
}

func TestParseTags(t *testing.T) {
	v := New()
	testCases := []struct {
		name     string
		input    string
		expected []Tag
	}{
		{
			"AllOK", "nonzero=true,regexp=[a-zA-Z]+", []Tag{
				Tag{Name: "nonzero", Fn: ValidateNonzero, Arg: "true"},
				Tag{Name: "regexp", Fn: ValidateRegexp, Arg: "[a-zA-Z]+"},
			},
		},
	}
	for _, tc := range testCases {
		got, _ := v.parseTags(tc.input)
		if !cmp(tc.expected, got) {
			t.Errorf("Both slices should be equal, got: %+v, expected: %+v", got, tc.expected)
		}
	}
	rawTags := "blah=omg"
	_, err := v.parseTags(rawTags)
	if err != ErrUnknownTag {
		t.Errorf("Should raise ErrUnknownTag error.")
	}
}

func TestValidate(t *testing.T) {
	v := New()
	testCases := []struct {
		name     string
		input    interface{}
		expected error
	}{
		{
			"AgeZeroValue", struct {
				Name string `validate:"regexp=[a-zA-Z]+"`
				Age  string `validate:"nonzero=true"`
			}{Name: "Admin"},
			ErrZeroValue,
		},
		{
			"RegexpNotMatch", struct {
				Name string `validate:"regexp=[a-zA-Z]+"`
				Age  string `validate:"nonzero=true"`
			}{Name: "123", Age: "12"},
			ErrRegexp,
		},
		{
			"AllOK", struct {
				Name string
				Age  string
			}{Name: "123", Age: "12"},
			nil,
		},
		{
			"InvalidEmail", struct {
				Email string `validate:"regexp=^\\S+@\\S+\\.\\S+$,nonzero=true"`
			}{Email: "john.doe"},
			ErrRegexp,
		},
		{
			"EmptyEmail", struct {
				Email string `validate:"nonzero=true"`
			}{Email: ""},
			ErrZeroValue,
		},
		{
			"FieldNotSupported", struct {
				Age int `validate:"nonzero=true"`
			}{Age: 100},
			ErrNotSupported,
		},
		{
			"EmptyField", struct {
				Name string `validate:"nonzero=true"`
			}{Name: ""},
			ErrZeroValue,
		},
		{
			"MaxLenExceeded",
			struct {
				Name string `validate:"maxlen=3"`
			}{Name: "1234"},
			ErrMaxLen,
		},
		{
			"MaxLenWrongArg",
			struct {
				Name string `validate:"maxlen=abc"`
			}{Name: "1234"},
			ErrInvalidArg,
		},
		{
			"MaxLenOk",
			struct {
				Name string `validate:"maxlen=3"`
			}{Name: "123"},
			nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := v.Validate(tc.input)
			if err != tc.expected {
				t.Errorf("Expected error: %v, got: %v", tc.expected, err)
			}
		})
	}
}
