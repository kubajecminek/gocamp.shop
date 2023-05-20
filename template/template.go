package template

import (
	"embed"
	"html/template"
	"net/http"
)

const (
	ErrServer = "Oops, tohle je chyba na straně serveru. Prosím kontaktujte mě: esprit.pristani0z [zavináč] icloud.com"
)

//go:embed static public
var content embed.FS

var funcs template.FuncMap = template.FuncMap{"mul": MulFunc, "add": AddFunc}

type Template struct {
	Content embed.FS
	Funcs   template.FuncMap
}

func New() *Template {
	return &Template{Content: content, Funcs: funcs}
}

func (t *Template) newTemplate(name string, pattern string) (*template.Template, error) {
	tmpl, err := template.New(name).Funcs(t.Funcs).ParseFS(t.Content, pattern, "static/layout.html")
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func (t *Template) Render(w http.ResponseWriter, data any, pattern string, statusCode int) {
	tmpl, err := t.newTemplate("tmpl", pattern)
	if err != nil {
		http.Error(w, ErrServer, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, ErrServer, http.StatusInternalServerError)
		return
	}
}
