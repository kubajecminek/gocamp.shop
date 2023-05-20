package web

import (
	"io/fs"
	"net/http"
)

func (w *Web) PublicFileServer(rw http.ResponseWriter, r *http.Request, sd SessionData) *Error {
	publicDir, err := fs.Sub(w.Template.Content, "public")
	if err != nil {
		return &Error{Err: err, Message: ErrMsg, Code: http.StatusInternalServerError}
	}
	http.StripPrefix("/public/", http.FileServer(http.FS(publicDir))).ServeHTTP(rw, r)
	return nil
}
