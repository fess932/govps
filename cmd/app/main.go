package main

import (
	"github.com/go-chi/chi/v5"
	"govps/pkg/server"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	r := chi.NewRouter()
	srv := server.NewDelivery(server.NewUsecase("/var/run/libvirt/libvirt-sock"))

	fsrv := http.FileServer(http.Dir("./web/"))
	r.Handle("/static/*", http.StripPrefix("/static/", fsrv))

	r.Get("/", srv.Home)

	r.Get("/vm/{id}", srv.VMInfo)
	r.Delete("/vm/{id}", srv.Delete)

	r.Get("/add", srv.AddVMPage)
	r.Post("/add", srv.Add)

	log.Println("server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
