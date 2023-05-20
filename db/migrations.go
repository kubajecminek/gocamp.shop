package db

import (
	"cloud.google.com/go/bigquery"
	"google.golang.org/api/googleapi"
	"net/http"
)

func (db Db) CreateOrdersTable() error {
	t := db.Client.Dataset("web").Table("orders")
	billingSchema := bigquery.Schema{
		{Name: "name", Type: bigquery.StringFieldType, Required: true},
		{Name: "street", Type: bigquery.StringFieldType, Required: true},
		{Name: "city", Type: bigquery.StringFieldType, Required: true},
		{Name: "zip_code", Type: bigquery.StringFieldType, Required: true},
		{Name: "id", Type: bigquery.StringFieldType},
		{Name: "tax_id", Type: bigquery.StringFieldType},
	}

	participantsSchema := bigquery.Schema{
		{Name: "item_id", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "name", Type: bigquery.StringFieldType, Required: true},
		{Name: "id", Type: bigquery.StringFieldType, Required: true},
		{Name: "club", Type: bigquery.StringFieldType},
		{Name: "weight_category", Type: bigquery.StringFieldType},
		{Name: "belt", Type: bigquery.StringFieldType},
		{Name: "shirt_size", Type: bigquery.StringFieldType, Required: true},
		{Name: "extra_requirements", Type: bigquery.StringFieldType},
	}

	checkoutSchema := bigquery.Schema{
		{Name: "name", Type: bigquery.StringFieldType, Required: true},
		{Name: "mobile_number", Type: bigquery.StringFieldType, Required: true},
		{Name: "email", Type: bigquery.StringFieldType, Required: true},
		{Name: "billing", Type: bigquery.RecordFieldType, Required: true, Repeated: false, Schema: billingSchema},
		{Name: "participants", Type: bigquery.RecordFieldType, Repeated: true, Schema: participantsSchema},
		{Name: "note", Type: bigquery.StringFieldType},
	}

	itemSchema := bigquery.Schema{
		{Name: "id", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "name", Type: bigquery.StringFieldType, Required: true},
		{Name: "description", Type: bigquery.StringFieldType, Required: true},
		{Name: "category", Type: bigquery.StringFieldType, Required: true},
		{Name: "price", Type: bigquery.IntegerFieldType, Required: true},
	}

	cartSchema := bigquery.Schema{
		{Name: "item", Type: bigquery.RecordFieldType, Required: true, Repeated: false, Schema: itemSchema},
		{Name: "quantity", Type: bigquery.IntegerFieldType, Required: true},
	}

	tableSchema := bigquery.Schema{
		{Name: "created_at", Type: bigquery.TimestampFieldType, Required: true},
		{Name: "id", Type: bigquery.StringFieldType, Required: true},
		{Name: "checkout", Type: bigquery.RecordFieldType, Required: true, Repeated: false, Schema: checkoutSchema},
		{Name: "cart", Type: bigquery.RecordFieldType, Repeated: true, Schema: cartSchema},
		{Name: "variable_symbol", Type: bigquery.StringFieldType, Required: true},
		{Name: "status", Type: bigquery.StringFieldType, Required: true},
	}

	metaData := &bigquery.TableMetadata{
		Schema: tableSchema,
	}

	// Create the table if not already exists
	if _, err := t.Metadata(db.Ctx); err != nil {
		if e, ok := err.(*googleapi.Error); ok {
			if e.Code == http.StatusNotFound {
				if err = t.Create(db.Ctx, metaData); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
