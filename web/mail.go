package web

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"net/mail"
	"net/smtp"
	"os"
	"text/tabwriter"
	"text/template"

	mytemplate "gocamp.shop/template"
)

func (w *Web) SendConfirmationEmail(td TemplateData) error {
	passwd, ok := os.LookupEnv("SMTP_PASSWORD")
	if !ok {
		return errors.New("env variable 'SMTP_PASSWORD' not set")
	}
	auth := smtp.PlainAuth("", w.Shop.Email.SenderAddress, passwd, w.Shop.Email.Host)

	to := []string{td.Order.Checkout.Email}
	msg, err := w.msgBody(td)
	if err != nil {
		return err
	}

	fullHost := fmt.Sprintf("%s:%d", w.Shop.Email.Host, w.Shop.Email.Port)
	if err := smtp.SendMail(fullHost, auth, w.Shop.Email.SenderAddress, to, msg); err != nil {
		return err
	}
	return nil
}

func (w *Web) msgBody(td TemplateData) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	t, err := template.New("email").Funcs(template.FuncMap{"mul": mytemplate.MulFunc}).Parse(confirmationEmailTemplate)
	if err != nil {
		return nil, err
	}
	tw := tabwriter.NewWriter(buf, 8, 8, 8, ' ', 0)
	if err := t.Execute(tw, td); err != nil {
		return nil, err
	}
	if err := tw.Flush(); err != nil {
		return nil, err
	}

	from := mail.Address{Name: w.Shop.Email.SenderName, Address: w.Shop.Email.SenderAddress}
	to := mail.Address{Name: "", Address: td.Order.Checkout.Email}
	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = mime.QEncoding.Encode("UTF-8", fmt.Sprintf("Děkujeme za objednávku č. %v", td.Order.VS()))
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString(buf.Bytes())
	return []byte(message), nil
}

var confirmationEmailTemplate = `
Dobrý den,
děkujeme za Vaší objednávku v našem obchodě. Vaše objednávka č. {{ ($.Order.VS) }} byla v pořádku přijata.

Platební údaje:
Částka: {{ ($.Order.TotalPrice) }} Kč
Variabilní symbol: {{ ($.Order.VS) }}
Číslo účtu: {{ .BankAccount }}

Kontaktní údaje:
Jméno a příjmení: {{ .Order.Checkout.Name }}
Mobil: {{ .Order.Checkout.MobileNumber}}
E-mail: {{ .Order.Checkout.Email }}

Fakturační údaje:
Jméno a přijmení (název firmy): {{ .Order.Checkout.Billing.Name }}
Ulice a číslo popisné: {{ .Order.Checkout.Billing.Street }}
Město: {{ .Order.Checkout.Billing.City }}
PSČ: {{ .Order.Checkout.Billing.ZipCode }}
{{- if .Order.Checkout.Billing.ID }}
IČO: {{ .Order.Checkout.Billing.ID }}
{{- end }}
{{- if .Order.Checkout.Billing.TaxID }}
DIČ: {{ .Order.Checkout.Billing.TaxID }}
{{- end }}

Košík:
{{- range .Order.Cart }}
{{ .Item.Description }} 	{{ .Quantity }} ks 	{{ .Item.Price }} Kč/ks 	{{ (mul .Item.Price .Quantity) }} Kč
{{ end }}


{{- if .Order.Checkout.Participants }}
Účastníci:
{{- end }}
{{ range $value := .Order.Checkout.Participants }}
Kemp: {{ ($.Order.ItemDescByID $value.ItemID) }}, {{ $value.OrdinalPos }}. účastník
Jméno a příjmení: {{ $value.Name }}
Rodné číslo: {{ $value.ID }}
{{- if $value.Club }}
Klub: {{ $value.Club }}
{{- end }}
{{- if $value.WeightCategory }}
Váhová kategorie: {{ $value.WeightCategory }}
{{- end }}
{{- if $value.Belt }}
Dosažený technický stupeň: {{ $value.Belt }}
{{- end }}
Účastnické tričko: {{ $value.ShirtSize }}
{{- if $value.ExtraRequirements }}
Speciální požadavky: {{ $value.ExtraRequirements }}
{{- end }}
{{ end }}

{{- if .Order.Checkout.Note }}
Poznámka:
{{ .Order.Checkout.Note }}
{{- end }}`
