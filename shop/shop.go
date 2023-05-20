package shop

import (
	"embed"
	"encoding/json"

	"gocamp.shop/db/models"
)

type Shop struct {
	// ProjectID is the ID of the Google Cloud project.
	ProjectID string `json:"projectId"`
	// Items is a list of items that can be purchased.
	Items []models.Item `json:"items"`
	// Name is the name of the shop.
	Name string `json:"name"`
	// BankAccount is the bank account number of the shop.
	BankAccount string `json:bankAccount`
	// Email is the email configuration of the shop.
	Email Email `json:"email"`
	// Stylesheet is the CSS stylesheet used by the shop.
	Stylesheet string `json:"stylesheet"`
	// FavIcon is the favicon used by the shop.
	FavIcon string `json:"favIcon"`
}

type Email struct {
	// Host is the SMTP host used to send emails.
	Host string `json:"host"`
	// Port is the SMTP port used to send emails.
	Port int `json:"port"`
	// SenderAddress is the email address used to send emails.
	SenderAddress string `json:"senderAddress"`
	// SenderName is the name used to send emails.
	SenderName string `json:"senderName"`
}

//go:embed config.json
var file embed.FS

func NewShop() (Shop, error) {
	var shop Shop
	configFile, err := file.ReadFile("config.json")
	if err != nil {
		return shop, err
	}
	if err := json.Unmarshal(configFile, &shop); err != nil {
		return shop, err
	}
	return shop, nil
}
