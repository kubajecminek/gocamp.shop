package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (w *Web) CartAdd(rw http.ResponseWriter, r *http.Request, sd SessionData) *Error {
	err := r.ParseForm()
	if err != nil {
		return &Error{Err: err, Message: ErrParseFormMsg, Code: http.StatusBadRequest}
	}
	for k, v := range r.Form {
		quantity, err := strconv.Atoi(v[0])
		if err != nil {
			return &Error{Err: err, Message: ErrParseFormMsg, Code: http.StatusBadRequest}
		}
		item, ok := w.ItemByName(k)
		if !ok {
			return &Error{Err: fmt.Errorf("web: item not found in sortiment %s", k), Message: ErrParseFormMsg, Code: http.StatusBadRequest}
		}
		sd.Order.Cart.Add(item, quantity)
	}
	w.SaveSessionData(r, sd)
	http.Redirect(rw, r, "/", http.StatusSeeOther)
	return nil
}

func (w *Web) AddToCart(rw http.ResponseWriter, r *http.Request, sd SessionData) *Error {
	op := chi.URLParam(r, "op")

	id, err := strconv.Atoi(chi.URLParam(r, "itemSlug"))
	if err != nil {
		return &Error{Err: err, Message: ErrParamMsg, Code: http.StatusBadRequest}
	}
	item, ok := w.ItemByID(id)
	if !ok {
		return &Error{Err: fmt.Errorf("web: item not found in sortiment %d", id), Message: ErrItemNotFoundMsg, Code: http.StatusBadRequest}
	}

	switch op {
	case "add":
		sd.Order.Cart.Add(item, 1)
	case "remove":
		sd.Order.Cart.Add(item, -1)
	}

	w.SaveSessionData(r, sd)

	http.Redirect(rw, r, "/", http.StatusSeeOther)
	return nil
}
