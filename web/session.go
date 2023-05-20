package web

import (
	"net/http"

	"gocamp.shop/db/models"
)

const (
	sessionDataKey = "GoCampShop"
)

var (
	ErrSessionDataNotValid = NewWebError(nil, ErrMsg, http.StatusInternalServerError)
	ErrSessionDataNotFound = NewWebError(nil, "Session data not found", http.StatusOK)
)

type SessionData struct {
	Order     models.Order
	LastOrder models.Order
}

func (sd *SessionData) OrderFinished() {
	sd.LastOrder = sd.Order
	sd.Order = models.NewOrder()
}

func (w *Web) SaveSessionData(r *http.Request, newData SessionData) {
	w.SessionManager.Put(r.Context(), sessionDataKey, newData)
}

func (w *Web) GetSessionData(r *http.Request) (SessionData, *Error) {
	if !w.SessionManager.Exists(r.Context(), sessionDataKey) {
		return SessionData{
			Order: models.NewOrder(),
		}, ErrSessionDataNotFound
	}
	sd, ok := w.SessionManager.Get(r.Context(), sessionDataKey).(SessionData)
	if !ok {
		return SessionData{
			Order: models.NewOrder(),
		}, ErrSessionDataNotValid
	}
	return sd, nil
}

func (w *Web) FlushSession(r *http.Request) *Error {
	err := w.SessionManager.Clear(r.Context())
	if err != nil {
		return &Error{Err: err, Message: ErrMsg, Code: http.StatusInternalServerError}
	}
	return nil
}
