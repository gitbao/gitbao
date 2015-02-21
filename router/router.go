package router

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gitbao/gitbao/model"
	"github.com/jinzhu/gorm"
)

var destinations map[string]string

func init() {
	err := populateDestinations()
	if err != nil {
		panic(err)
	}
}

type Router struct{}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	host := req.Host
	// host = strings.TrimSpace(host)
	//Figure out if a subdomain exists in the host given.
	host_parts := strings.Split(host, ".")
	if len(host_parts) > 2 {
		//The subdomain exists, we store it as the first element
		//in a new array
		subdomain := host_parts[0]
		// subdomain = "ba67234b79784c75cfd9-1"
		destination, err := GetDestinaton(subdomain)
		if err == gorm.RecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Not Found"))
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		_ = destination
		// Log if requested

		// We'll want to use a new client for every request.
		client := &http.Client{}

		// Tweak the request as appropriate:
		//	RequestURI may not be sent to client
		//	URL.Scheme must be lower-case
		req.RequestURI = ""
		req.URL.Scheme = "http"
		// req.URL = destination
		req.URL.Host = destination
		// And proxy
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		io.Copy(w, resp.Body)
		// resp.Write(w)
		return
	}
	w.WriteHeader(http.StatusNotFound)

}

func GetDestinaton(subdomain string) (
	destination string, err error) {

	val, exists := destinations[subdomain]
	if exists == false {
		var location model.Location
		query := model.DB.Where("subdomain = ?", subdomain).
			Find(&location)
		if err = query.Error; err != nil {
			return
		}
		destination = location.Destination
		return
	}

	destination = val
	return
}

func populateDestinations() (err error) {
	var locations []model.Location
	query := model.DB.Find(&locations)
	if query.Error != nil {
		err = query.Error
		return
	}
	destinations = make(map[string]string)
	for _, value := range locations {
		destinations[value.Subdomain] = value.Destination
	}
	return
}
