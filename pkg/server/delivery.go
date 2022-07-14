package server

import (
	"encoding/hex"
	"fmt"
	"github.com/digitalocean/go-libvirt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log"
	"net/http"
	"strconv"
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
				switch val := v.(type) {
				case libvirt.UUID:
					return hex.EncodeToString(val[:])
				case [32]int8:
					return fmt.Sprintf("%v", val)
				case string:
					return val
				default:
					return fmt.Sprintf("%v", v)
				}
			},
			"memoryToString": func(v uint64) string {
				return fmt.Sprintf("%v МБ", v/1024)
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

func (d *Delivery) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		HTTPError(w, http.StatusBadRequest, fmt.Errorf("wrong id %v: %w", id, err))
		return
	}

	if err = d.uc.Delete(int32(id)); err != nil {
		HTTPError(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
