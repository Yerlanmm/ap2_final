package handler

import (
	"html/template"
	"net/http"
)

func PaymentPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/payment.html"))
	tmpl.Execute(w, nil)
}
