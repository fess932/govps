package server

import (
	"html/template"
	"log"
	"net/http"
)

func HTTPError(w http.ResponseWriter, status int, err error) {
	log.Printf("HTTPError: %v", err)
	http.Error(w, http.StatusText(status), status)
}

func NewDelivery(uc *Usecase) *Delivery {
	return &Delivery{uc}
}

type Delivery struct {
	uc *Usecase
}

func (d *Delivery) Home(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./web/template/home.html")
	if err != nil {
		HTTPError(w, http.StatusInternalServerError, err)
		return
	}

	h, err := d.uc.Get()
	if err != nil {
		HTTPError(w, http.StatusInternalServerError, err)
		return
	}

	if err = t.Execute(w, h); err != nil {
		HTTPError(w, http.StatusInternalServerError, err)
		return
	}
}
