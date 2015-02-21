package router

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
	r := &Router{}
	r.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Error(fmt.Errorf("Wrong status code"))
	}
	if string(w.Body.Bytes()) != "404 page not found\n" {
		t.Error(fmt.Errorf("Wrong response body"))
	}
}

func TestPopulateDestinations(t *testing.T) {
	populateDestinations()

	if len(destinations) != 4 {
		t.Error("Wrong number of destinations")
	}
}
func TestGetDestiantionFromMap(t *testing.T) {

	// test that all the destinations are
	// populated
	for key, value := range destinations {
		destination, err := GetDestinaton(key)
		if err != nil {
			t.Error(err)
		}
		if value != destination {
			t.Error(
				fmt.Errorf(
					"Destination %s incorrectly populated",
					key,
				),
			)
		}
	}

}
func TestGetDestiantionFromDB(t *testing.T) {

	// add a new location to the database
	// and then test for it.
	location := model.Location{
		Subdomain:   "unpopulated",
		Destination: "4.4.4.4",
	}
	query := model.DB.Create(&location)
	if query.Error != nil {
		panic(query.Error)
	}
	destination, err := GetDestinaton("unpopulated")
	if err != nil {
		t.Error(err)
	}
	if "4.4.4.4" != destination {
		t.Error(
			fmt.Errorf(
				"Database-based location is incorrectly retrieved",
			),
		)
	}
}

func BenchmarkGetDestinationFromDB(b *testing.B) {
	// run the Fib function b.N times
	var t *testing.T
	for n := 0; n < b.N; n++ {
		TestGetDestiantionFromDB(t)
	}
}

func BenchmarkGetDestinationFromMap(b *testing.B) {
	// run the Fib function b.N times
	var t *testing.T
	for n := 0; n < b.N; n++ {
		TestGetDestiantionFromMap(t)
	}
}
