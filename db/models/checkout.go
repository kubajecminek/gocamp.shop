package models

import (
	"fmt"
	"net/url"

	"gocamp.shop/validator"
)

type Checkout struct {
	Name         string        `id:"name" validate:"nonzero=true,maxlen=30" bigquery:"name"`
	MobileNumber string        `id:"mnumber" validate:"nonzero=true,maxlen=30" bigquery:"mobile_number"`
	Email        string        `id:"email" validate:"nonzero=true,regexp=^\\S+@\\S+\\.\\S+$,maxlen=80" bigquery:"email"`
	Billing      Billing       `bigquery:"billing"`
	Participants []Participant `bigquery:"participants"`
	Note         string        `id:"note" validate:"maxlen=200" bigquery:"note"`
}

func (c *Checkout) setBilling(billing Billing) {
	c.Billing = billing
}

func (c *Checkout) setParticipants(participants []Participant) {
	c.Participants = participants
}

func decodeCheckout(form url.Values) Checkout {
	return Checkout{
		Name:         form.Get("name"),
		MobileNumber: form.Get("mnumber"),
		Email:        form.Get("email"),
		Billing:      Billing{},
		Participants: []Participant{},
		Note:         form.Get("note"),
	}
}

func (c Checkout) Validate(v validator.Validator, minParticipants int) error {
	if err := v.Validate(c.Billing); err != nil {
		return err
	}
	if len(c.Participants) < minParticipants {
		return fmt.Errorf("models: invalid number of participants, have %d, want %d", len(c.Participants), minParticipants)
	}
	for _, p := range c.Participants {
		if err := v.Validate(p); err != nil {
			return err
		}
	}
	if err := v.Validate(c); err != nil {
		return err
	}
	return nil
}

func (c Checkout) IsValid(v validator.Validator, minParticipants int) bool {
	return c.Validate(v, minParticipants) == nil
}

func (c Checkout) IsEmpty() bool {
	return c.Name == "" && c.MobileNumber == "" && c.Email == "" && c.Billing.IsEmpty() && len(c.Participants) == 0 && c.Note == ""
}
