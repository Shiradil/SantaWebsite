package handlers

import (
	"html/template"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) error {
	t, err := template.ParseFiles("ui/templates/" + tmpl)
	if err != nil {
		return err
	}

	err = t.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}
