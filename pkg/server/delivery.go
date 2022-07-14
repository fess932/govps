package server

import (
	"encoding/hex"
	"github.com/digitalocean/go-libvirt"
	"html/template"
	"log"
	"net/http"
)

func HTTPError(w http.ResponseWriter, status int, err error) {
	log.Printf("HTTPError: %v", err)
	http.Error(w, err.Error(), status)
}

func NewDelivery(uc *Usecase) *Delivery {
	return &Delivery{uc}
}

type Delivery struct {
	uc *Usecase
}

func (d *Delivery) Home(w http.ResponseWriter, r *http.Request) {
	t := template.New("home.html").
		Funcs(template.FuncMap{
			"toString": func(v interface{}) string {
				uuid, ok := v.(libvirt.UUID)
				if !ok {
					return "unknown type"
				}

				return hex.EncodeToString(uuid[:])
			},
		})

	t, err := t.ParseFiles("./web/template/home.html")
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

func (d *Delivery) Add(w http.ResponseWriter, r *http.Request) {
	if err := d.uc.Create(); err != nil {
		HTTPError(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
