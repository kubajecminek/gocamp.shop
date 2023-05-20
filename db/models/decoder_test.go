package models

import (
	"net/url"
	"reflect"
	"testing"
)

func TestDecode(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected Checkout
	}{
		{"simple", "pname-56-01=John", Checkout{Participants: []Participant{{Suffix: "-56-01", OrdinalPos: 56, ItemID: 1, Name: "John"}}}},
		{"twoParticipants", "pname-56-01=John&pname-56-02=Jane",
			Checkout{
				Participants: []Participant{
					{Suffix: "-56-01", OrdinalPos: 56, ItemID: 1, Name: "John"},
					{Suffix: "-56-02", OrdinalPos: 56, ItemID: 2, Name: "Jane"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			form, _ := url.ParseQuery(tc.input)
			got, err := Decode(form)
			if err != nil {
				t.Errorf("case: %s, got error %v", tc.name, err)
			}
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("case: %s, got %v, expected %v", tc.name, got, tc.expected)
			}
		})
	}
}
