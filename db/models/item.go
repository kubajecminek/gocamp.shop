package models

import "strings"

type Item struct {
	ID          int    `json:"id" bigquery:"id"`
	Name        string `json:"name" bigquery:"name"`
	Description string `json:"description" bigquery:"description"`
	Category    string `json:"category" bigquery:"category"`
	Price       int    `json:"price" bigquery:"price"`
	Img         string `json:"img" bigquery:"-"`
}

func (i Item) IsCamp() bool {
	return strings.Contains(i.Category, "camp")
}

func (i Item) IsBackprint() bool {
	return strings.Contains(i.Category, "backprint")
}
