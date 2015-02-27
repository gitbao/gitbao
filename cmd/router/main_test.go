package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gitbao/gitbao/model"
)

func init() {
	destinations := map[string]string{
		"localhost": "127.128.1.1",
		"localhos":  "127.128.1.3",
		"localho":   "127.128.1.5",
		"localh":    "127.0.0.1:8000",
	}

	for key, value := range destinations {
		location := model.Location{
			Subdomain:   key,
			Destination: value,
		}
		query := model.DB.Create(&location)
		if query.Error != nil {
			panic(query.Error)
		}
	}
}

func TestDownloadHandler(t *testing.T) {
	go http.ListenAndServe(":8000", nil)
	req, err := http.NewRequest(
		"GET",
		"https://localh.gitbao.com/",
		nil,
	)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	r := Handler()
	r.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Error(fmt.Errorf("Wrong status code"))
	}
	if string(w.Body.Bytes()) != "404 page not found\n" {
		t.Error(fmt.Errorf("Wrong response body"))
	}
}
