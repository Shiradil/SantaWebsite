package handlers

import (
	"html/template"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) error {
	t, err := template.ParseFiles("ui/templates/" + tmpl)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
