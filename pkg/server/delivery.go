package server

import (
	"encoding/hex"
	"fmt"
	"github.com/digitalocean/go-libvirt"
	"github.com/go-chi/chi/v5"
	"govps/pkg"
	"html/template"
	"net/http"
)

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
		pkg.HTTPError(w, http.StatusInternalServerError, err)
		return
	}

	h, err := d.uc.Get()
	if err != nil {
		pkg.HTTPError(w, http.StatusInternalServerError, err)
		return
	}

	if err = t.Execute(w, h); err != nil {
		pkg.HTTPError(w, http.StatusInternalServerError, err)
		return
	}
}

func (u *Usecase) CreateVMPage(w http.ResponseWriter, r *http.Request) {
	t := template.New("create.html").
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

	t, err := t.ParseFiles("./web/template/create.html")
	if err != nil {
		pkg.HTTPError(w, http.StatusInternalServerError, err)
		return
	}

	if err = t.Execute(w, nil); err != nil {
		pkg.HTTPError(w, http.StatusInternalServerError, err)
		return
	}
}

func (d *Delivery) Add(w http.ResponseWriter, r *http.Request) {
	if err := d.uc.Create(); err != nil {
		pkg.HTTPError(w, http.StatusInternalServerError, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func strToUUID(str string) (libvirt.UUID, error) {
	var uuid libvirt.UUID

	buf, err := hex.DecodeString(str)
	if err != nil {
		return uuid, fmt.Errorf("wrong id: %w", err)
	}

	copy(uuid[:], buf)

	return uuid, nil
}

func (d *Delivery) Delete(w http.ResponseWriter, r *http.Request) {
	uuid, err := strToUUID(chi.URLParam(r, "id"))
	if err != nil {
		pkg.HTTPError(w, http.StatusBadRequest, fmt.Errorf("wrong id: %w", err))
		return
	}

	if err = d.uc.Delete(uuid); err != nil {
		pkg.HTTPError(w, http.StatusInternalServerError, err)
		return
	}
}

func (d *Delivery) AddVMPage(w http.ResponseWriter, r *http.Request) {
	t := template.New("create.html")

	t, err := t.ParseFiles("./web/template/create.html")
	if err != nil {
		pkg.HTTPError(w, http.StatusInternalServerError, err)
		return
	}

	h, err := d.uc.Get()
	if err != nil {
		pkg.HTTPError(w, http.StatusInternalServerError, err)
		return
	}

	if err = t.Execute(w, h); err != nil {
		pkg.HTTPError(w, http.StatusInternalServerError, err)
		return
	}
}

func (d *Delivery) VMInfo(w http.ResponseWriter, r *http.Request) {
	id, err := strToUUID(chi.URLParam(r, "id"))
	if err != nil {
		pkg.HTTPError(w, http.StatusBadRequest, fmt.Errorf("wrong id %v: %w", id, err))
		return
	}

	t := template.New("vm.html").Funcs(template.FuncMap{
		"state": func(v uint8) string {
			switch v {
			case 1:
				return "Запущен"
			case 2:
				return "Заблокирован"
			case 3:
				return "Пауза"
			case 4:
				return "Выключен"
			case 6:
				return "Крашнулся"
			default:
				return "unknown"
			}
		},
		"memoryToString": func(v uint64) string {
			return fmt.Sprintf("%v МБ", v/1024)
		},
	})

	t, err = t.ParseFiles("./web/template/vm.html")
	if err != nil {
		pkg.HTTPError(w, http.StatusInternalServerError, err)
		return
	}

	info, err := d.uc.GetVMInfo(id)
	if err != nil {
		pkg.HTTPError(w, http.StatusInternalServerError, err)
		return
	}

	if err = t.Execute(w, pkg.JSON{"id": id, "info": info}); err != nil {
		pkg.HTTPError(w, http.StatusInternalServerError, err)
		return
	}
}
