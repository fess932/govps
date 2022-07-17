package pkg

import (
	"log"
	"net/http"
)

type JSON map[string]interface{}

func HTTPError(w http.ResponseWriter, status int, err error) {
	log.Printf("HTTPError: %v", err)
	http.Error(w, err.Error(), status)
}
