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
	r.Get("/add", srv.Add)
	r.Get("/delete/{id}", srv.Delete)

	log.Println("server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
