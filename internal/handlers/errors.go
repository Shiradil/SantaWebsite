package handlers

import (
	"html/template"
	"net/http"
)

type errorss struct {
	ErrorCode int
	ErrorMsg  string
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, errCode int, msg string) {
	t, err := template.ParseFiles("ui/templates/Error.html")
	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		ErrorHandler(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}
	Errors := errorss{
		ErrorCode: errCode,
		ErrorMsg:  msg,
	}
	// w.WriteHeader(Errors.ErrorCode)
	t.Execute(w, Errors)
}
