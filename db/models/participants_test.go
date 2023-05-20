package models

import (
	"net/url"
	"reflect"
	"testing"
)

func TestAttrSuffixes(t *testing.T) {
	// Test data
	testCases := []struct {
		name string
		in   string
		out  []string
	}{
		{name: "simple", in: "pname-0-1=A", out: []string{"-0-1"}},
		{name: "twoParticipants", in: "pname-0-1=A&pid-1-1=123", out: []string{"-0-1", "-1-1"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the code we are testing
			form, _ := url.ParseQuery(tc.in)
			participants, _ := attrSuffixes(form)
			// Check the result
			if !reflect.DeepEqual(participants, tc.out) {
				t.Errorf("expected %v, actual %v", tc.out, participants)
			}
		})
	}

	f2 := url.Values{"pname-abc-def": []string{"A"}}
	_, err := attrSuffixes(f2)
	if err.Error() != "models: invalid attribute pname-abc-def" {
		t.Errorf("expected %v, actual %v", "models: invalid attribute pname-abc-def", err.Error())
	}
}

func TestAttrMatch(t *testing.T) {
	testCases := []struct {
		name     string
		attr     string
		slice    []string
		expected bool
	}{
		{"oneTag", "id-0-1", []string{"id"}, true},
		{"emptySlice", "id-0-1", []string{}, false},
		{"emptyAttr", "", []string{"id"}, false},
		{"twoMatches", "pname-035-1", []string{"pname", "id", "pname"}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := attrMatch(tc.attr, tc.slice)
			if got != tc.expected {
				t.Errorf("case: %s, got %v, expected %v", tc.name, got, tc.expected)
			}
		})
	}
}

func TestDecodeParticipants(t *testing.T) {
	// Test data
	testCases := []struct {
		name string
		in   string
		out  []Participant
	}{
		{
			name: "simple",
			in:   "pname-0-1=A&pid-0-1=123",
			out:  []Participant{Participant{Suffix: "-0-1", OrdinalPos: 0, ItemID: 1, Name: "A", ID: "123"}}},
		{
			name: "twoParticipants",
			in:   "pname-0-1=A&pid-0-1=123&pname-1-1=B&pid-1-1=456",
			out: []Participant{
				Participant{Suffix: "-0-1", OrdinalPos: 0, ItemID: 1, Name: "A", ID: "123"},
				Participant{Suffix: "-1-1", OrdinalPos: 1, ItemID: 1, Name: "B", ID: "456"}}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the code we are testing
			form, _ := url.ParseQuery(tc.in)
			participants, _ := decodeParticipants(form)
			// Check the result
			if !reflect.DeepEqual(participants, tc.out) {
				t.Errorf("expected %v, actual %v", tc.out, participants)
			}
		})
	}
}
