package db

import (
	"context"

	"cloud.google.com/go/bigquery"
	"gocamp.shop/db/models"
)

type Db struct {
	Client         *bigquery.Client
	Ctx            context.Context
	OrdersInserter *bigquery.Inserter
}

func NewDB(ctx context.Context, client *bigquery.Client) Db {
	return Db{
		Client:         client,
		Ctx:            ctx,
		OrdersInserter: client.Dataset("web").Table("orders").Inserter(),
	}
}

func (db Db) InsertOrder(order models.Order) error {
	if err := db.OrdersInserter.Put(db.Ctx, order); err != nil {
		return err
	}
	return nil
}
