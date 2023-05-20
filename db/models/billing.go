package models

import (
	"net/url"
)

type Billing struct {
	Name    string `id:"bname" validate:"nonzero=true,maxlen=30" bigquery:"name"`
	Street  string `id:"bstreet" validate:"nonzero=true,maxlen=80" bigquery:"street"`
	City    string `id:"bcity" validate:"nonzero=true,maxlen=30" bigquery:"city"`
	ZipCode string `id:"bzipcode" validate:"nonzero=true,maxlen=10" bigquery:"zip_code"`
	ID      string `id:"binum" validate:"maxlen=20" bigquery:"id"`
	TaxID   string `id:"btnum" validate:"maxlen=20" bigquery:"tax_id"`
}

func decodeBilling(form url.Values) Billing {
	return Billing{
		Name:    form.Get("bname"),
		Street:  form.Get("bstreet"),
		City:    form.Get("bcity"),
		ZipCode: form.Get("bzipcode"),
		ID:      form.Get("binum"),
		TaxID:   form.Get("btnum"),
	}
}

func (b Billing) IsEmpty() bool {
	return b.Name == "" && b.Street == "" && b.City == "" && b.ZipCode == "" && b.ID == "" && b.TaxID == ""
}
