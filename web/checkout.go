package web

import (
	"errors"
	"net/http"

	"gocamp.shop/db/models"
)

func (w *Web) CheckoutForm(rw http.ResponseWriter, r *http.Request, sd SessionData) *Error {
	if sd.Order.Cart.IsEmpty() {
		http.Redirect(rw, r, "/", http.StatusFound)
		return nil
	}

	w.Template.Render(rw, newD(sd.Order, w.Shop), "static/checkout-form.html", 200)
	return nil
}

func (w *Web) Checkout(rw http.ResponseWriter, r *http.Request, sd SessionData) *Error {
	if sd.Order.Cart.IsEmpty() {
		return &Error{Err: errors.New("web: cart is empty"), Message: ErrCartEmptyMsg, Code: http.StatusBadRequest}
	}
	if err := r.ParseForm(); err != nil {
		return &Error{Err: err, Message: ErrParseFormMsg, Code: http.StatusBadRequest}
	}
	c, err := models.Decode(r.Form)
	if err != nil {
		return &Error{Err: err, Message: ErrCheckoutMsg, Code: http.StatusBadRequest}
	}
	if err = c.Validate(w.Validator, sd.Order.Cart.NumCamps()); err != nil {
		return &Error{Err: err, Message: ErrCheckoutMsg, Code: http.StatusBadRequest}
	}
	sd.Order.Checkout = c
	w.SaveSessionData(r, sd)

	w.Template.Render(rw, newD(sd.Order, w.Shop), "static/checkout.html", 200)
	return nil
}
