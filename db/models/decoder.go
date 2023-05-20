package models

import (
	"net/url"
)

// Decode decodes the form into a Checkout.
func Decode(form url.Values) (Checkout, error) {
	billing := decodeBilling(form)

	participants, err := decodeParticipants(form)
	if err != nil {
		return Checkout{}, err
	}

	checkout := decodeCheckout(form)
	checkout.setBilling(billing)
	checkout.setParticipants(participants)

	return checkout, nil
}
