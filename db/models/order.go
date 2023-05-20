package models

import (
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"time"
)

type Order struct {
	Created        time.Time `bigquery:"created_at"`
	ID             string    `bigquery:"id"`
	Checkout       Checkout  `bigquery:"checkout"`
	Cart           Cart      `bigquery:"cart"`
	VariableSymbol string    `bigquery:"variable_symbol"`
	Status         string    `bigquery:"status"`
}

func NewOrder() Order {
	uuid := uuid.New()

	return Order{
		Created:        time.Now(),
		ID:             uuid.String(),
		Cart:           NewCart(),
		VariableSymbol: strconv.FormatUint(uint64(uuid.ID()), 10),
		Status:         "in-progress",
	}
}

func (o Order) TotalPrice() int {
	var total int
	for _, item := range o.Cart {
		total += item.Quantity * item.Item.Price
	}
	return total
}

func (o Order) VS() string {
	return o.VariableSymbol
}

func (o Order) Completed() bool {
	return o.Status == "completed"
}

func (o *Order) MarkCompleted() {
	o.Status = "completed"
}

func (o Order) ItemDescByID(id int) (string, error) {
	for _, item := range o.Cart {
		if item.Item.ID == id {
			return item.Item.Description, nil
		}
	}
	return "", fmt.Errorf("item with id %d not found", id)
}
