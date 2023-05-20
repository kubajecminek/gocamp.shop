package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"gocamp.shop/web"
)

const (
	version = "0.0.4"
)

type webHandler func(rw http.ResponseWriter, r *http.Request, sd web.SessionData) *web.Error

func appHandler(w web.Web, fn webHandler) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		sd, err := w.GetSessionData(r)
		if err == web.ErrSessionDataNotValid {
			w.Logger.Printf("Method: %s, Error: %v", r.Method, err.Err.Error())
			w.Template.Render(rw, struct {
				Err   error
				Name  string
				Title string
			}{Err: err, Name: w.Shop.Name}, "static/error.html", err.Code)
		}
		err = fn(rw, r, sd)
		if err != nil {
			w.Logger.Printf("Method: %s, Error: %v", r.Method, err.Err.Error())
			w.Template.Render(rw, struct {
				Err   error
				Name  string
				Title string
			}{Err: err, Name: w.Shop.Name}, "static/error.html", err.Code)
		}
	}
}

func main() {
	// TODO: Add CSRF protection
	logger := log.New(log.Writer(), log.Prefix(), log.Flags())
	w, err := web.NewWeb(logger)
	if err != nil {
		logger.Fatalln(err)
	}
	err = w.Db.CreateOrdersTable()
	if err != nil {
		logger.Fatalln(err)
	}
	r := chi.NewRouter()
	r.Use(httprate.LimitByIP(1000, 20*time.Minute))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	gob.Register(web.SessionData{})

	// Index
	r.Get("/", appHandler(w, w.Index))

	// Cart
	r.Get("/cart/{op:(add)?(remove)?}/{itemSlug:[a-zA-Z0-9-]+}", appHandler(w, w.AddToCart))

	// Checkout
	r.Get("/checkout", appHandler(w, w.CheckoutForm))
	r.Post("/checkout", appHandler(w, w.Checkout))

	// Orders
	r.Post("/order", appHandler(w, w.OrderCreate))
	r.Get("/order-success", appHandler(w, w.OrderSuccess))

	// Misc
	r.Get("/public/*", appHandler(w, w.PublicFileServer))
	r.Get("/session/flush", appHandler(w, func(rw http.ResponseWriter, r *http.Request, sd web.SessionData) *web.Error {
		if err := w.FlushSession(r); err != nil {
			return err
		}
		http.Redirect(rw, r, "/", http.StatusSeeOther)
		return nil
	}))
	r.Get("/status", func(rw http.ResponseWriter, r *http.Request) { rw.Write([]byte("OK")) })
	r.Get("/version", func(rw http.ResponseWriter, r *http.Request) { rw.Write([]byte(version)) })

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
		logger.Printf("defaulting to port %s", port)
	}

	logger.Fatalln(http.ListenAndServe(":"+port, w.SessionManager.LoadAndSave(r)))
}
