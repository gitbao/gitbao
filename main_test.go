package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sqs/mux"
)

func TestDownloadHandler(t *testing.T) {
	req, err := http.NewRequest(
		"GET",
		"https://gist.github.com/maxmcd/ba67234b79784c75cfd9",
		nil,
	)
	if err != nil {
		t.Error(err)
	}

	m := mux.NewRouter()
	w := httptest.NewRecorder()
	m.HandleFunc("/{username}/{gist-id}", DownloadHandler)
	m.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Error(fmt.Errorf("Wrong response status code"))
	}
}

func TestMain(t *testing.T) {
	go main()
	resp, err := http.Get("http://localhost:8000/maxmcd/ba67234b79784c75cfd9")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%#v", resp)
	if resp.StatusCode != 200 {
		t.Error(fmt.Errorf("Wrong response status code"))
	}
}
