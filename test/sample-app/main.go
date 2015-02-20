package main

import (
	"net/http"

	"github.com/sqs/mux"
)

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}
