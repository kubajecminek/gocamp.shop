package web

import (
	"net/http"
)

func (w *Web) Index(rw http.ResponseWriter, r *http.Request, sd SessionData) *Error {
	w.Template.Render(rw, newD(sd.Order, w.Shop), "static/index.html", 200)
	return nil
}
