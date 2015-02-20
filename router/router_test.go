package router

import (
	"fmt"
	"testing"

	"github.com/gitbao/gitbao/model"
)

func init() {
	destinations := map[string]string{
		"localhost": "127.128.1.1",
		"localhos":  "127.128.1.3",
		"localho":   "127.128.1.5",
		"localh":    "127.128.1.6",
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
		destination, err := getDestinaton(key)
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
	destination, err := getDestinaton("unpopulated")
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
