package router

import "github.com/gitbao/gitbao/model"

var destinations map[string]string

func init() {
	err := populateDestinations()
	if err != nil {
		panic(err)
	}
}

func getDestinaton(subdomain string) (
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
