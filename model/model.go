package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var DB gorm.DB

type Location struct {
	Id          int64
	Subdomain   string
	Destination string
}

type Bao struct {
	Id         int64
	GistId     string
	Url        string
	GitPullUrl string
	Location   Location
	LocationId int64
}

func init() {

	var err error
	DB, err = gorm.Open("postgres", "dbname=gitbaotest sslmode=disable")
	if err != nil {
		panic(err)
	}

	DB.DropTableIfExists(&Location{})
	DB.DropTableIfExists(&Bao{})

	DB.AutoMigrate(&Location{}, &Bao{})
}
