package web

import (
	"context"
	"log"
	"time"

	"gocamp.shop/shop"

	"gocamp.shop/db"
	"gocamp.shop/db/models"

	"gocamp.shop/template"
	"gocamp.shop/validator"

	bq "cloud.google.com/go/bigquery"

	"github.com/alexedwards/scs/v2"
)

type Web struct {
	Shop           shop.Shop
	SessionManager *scs.SessionManager
	Template       *template.Template
	Validator      validator.Validator
	Db             db.Db
	Logger         *log.Logger
}

var SessionManager *scs.SessionManager

func NewWeb(logger *log.Logger) (Web, error) {
	// Initialize a new session manager and configure the session lifetime.
	SessionManager = scs.New()
	SessionManager.Lifetime = 24 * time.Hour

	v := validator.New()

	shop, err := shop.NewShop()
	if err != nil {
		return Web{}, err
	}

	ctx := context.Background()
	bqClient, err := bq.NewClient(ctx, shop.ProjectID)
	if err != nil {
		return Web{}, err
	}
	defer bqClient.Close()

	return Web{
		Shop:           shop,
		SessionManager: SessionManager,
		Template:       template.New(),
		Validator:      v,
		Db:             db.NewDB(ctx, bqClient),
		Logger:         logger,
	}, nil
}

func (w *Web) ItemByID(id int) (models.Item, bool) {
	for _, v := range w.Shop.Items {
		if id == v.ID {
			return v, true
		}
	}
	return models.Item{}, false
}

func (w *Web) ItemByName(name string) (models.Item, bool) {
	for _, v := range w.Shop.Items {
		if name == v.Name {
			return v, true
		}
	}
	return models.Item{}, false
}
