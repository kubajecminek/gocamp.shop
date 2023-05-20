# GoCamp Shop

<p align="center"><img src="/docs/preview.png"/></p>

`gocamp.shop` is a server-side application created for the purpose of Judo Heroes Camp. It is a simple e-commerce website where you can book a place at the judo camp and buy merchandise. The website features sessions, dynamic form generation, form validation, and more. It uses BigQuery as a backend database and requires no additional infrastructure. The entire configuration is done using a JSON file. An example configuration is shown below:

```
{
  "name": "Our shop",
  "bankAccount": "<bank-account>",
  "favIcon": "<favicon-url>",
  "projectID": "<bq-project-id>",
  "email": {
    "host": "<smtp-host>",
    "port": 587,
    "senderAddress": "john@shop.com",
    "senderName": "David from our shop"
  },
  "stylesheet": "<stylesheet-url>",
  "items": [
    {
      "id": 1,
      "name": "kemp-1",
      "description": "Winter camp | 1.–7. 1. 2020 (Prague)",
      "category": "winter_camp",
      "price": 5000,
      "img": "<img-url>"
    }
  ]
}
```

To set up the application, place this file into the `shop/` directory and run make build. Note that there are two environmental variables that must be set during runtime:
 - `SMTP_PASSWORD`: The password used for SMTP authentication (enables sending confirmation emails).
 - `GOOGLE_APPLICATION_CREDENTIALS`: The path to the service account with write access to the BigQuery table (`<project-id>.web.orders` by default).

Please make sure to set these environmental variables correctly before running the application.

## Project design
 - No javascript
 - Everything is handled server-side
 - Simple UI
 - Serverless DB
