package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type errorss struct {
	ErrorCode int
	ErrorMsg  string
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, errCode int, msg string) {
	t, err := template.ParseFiles("ui/templates/Error.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, strconv.Itoa(http.StatusInternalServerError)+" "+http.StatusText(http.StatusInternalServerError))
		log.Error(err.Error())
		return
	}
	Errors := errorss{
		ErrorCode: errCode,
		ErrorMsg:  msg,
	}
	t.Execute(w, Errors)
}
