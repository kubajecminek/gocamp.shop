package web

import (
	"gocamp.shop/db/models"
	"gocamp.shop/shop"
)

type TemplateData struct {
	Order               models.Order
	Sortiment           []models.Item
	Name                string
	Title               string
	Stylesheet          string
	SupportEmailAddress string
	BankAccount         string
	FavIcon             string
}

func (td *TemplateData) SetTitle(title string) {
	td.Title = title
}

func newD(order models.Order, shop shop.Shop) TemplateData {
	return TemplateData{
		Order:               order,
		Sortiment:           shop.Items,
		Name:                shop.Name,
		Title:               shop.Name,
		Stylesheet:          shop.Stylesheet,
		SupportEmailAddress: shop.Email.SenderAddress,
		BankAccount:         shop.BankAccount,
		FavIcon:             shop.FavIcon,
	}
}
