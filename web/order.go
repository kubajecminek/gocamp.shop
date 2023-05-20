package web

import (
	"errors"
	"net/http"
)

func (w *Web) OrderCreate(rw http.ResponseWriter, r *http.Request, sd SessionData) *Error {
	err := r.ParseForm()
	if err != nil {
		return &Error{Err: err, Message: ErrParseFormMsg, Code: http.StatusBadRequest}
	}
	agree := r.Form.Get("t-and-c")
	if !(agree == "agree") {
		return &Error{Err: errors.New("web: toc not agreed"), Message: ErrTOCMsg, Code: http.StatusBadRequest}
	}
	if sd.Order.Cart.IsEmpty() {
		return &Error{Err: errors.New("web: cart is empty"), Message: ErrCartEmptyMsg, Code: http.StatusBadRequest}
	}
	if err = sd.Order.Checkout.Validate(w.Validator, sd.Order.Cart.NumCamps()); err != nil {
		return &Error{Err: err, Message: ErrCheckoutMsg, Code: http.StatusBadRequest}
	}
	sd.Order.MarkCompleted()
	if err := w.Db.InsertOrder(sd.Order); err != nil {
		return &Error{Err: err, Message: ErrMsg, Code: http.StatusServiceUnavailable}
	}
	if err := w.SendConfirmationEmail(newD(sd.Order, w.Shop)); err != nil {
		return &Error{Err: err, Message: ErrMsg, Code: http.StatusInternalServerError}
	}
	sd.OrderFinished()
	w.SaveSessionData(r, sd)

	http.Redirect(rw, r, "/order-success", http.StatusFound)
	return nil
}

func (w *Web) OrderSuccess(rw http.ResponseWriter, r *http.Request, sd SessionData) *Error {
	if !sd.LastOrder.Completed() {
		http.Redirect(rw, r, "/", http.StatusFound)
		return nil
	}
	w.Template.Render(rw, newD(sd.LastOrder, w.Shop), "static/order.html", 200)
	return nil
}
