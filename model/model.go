package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var DB gorm.DB

type Config struct {
	BaoId   int64
	Port    int64
	EnvVars []EnvVar
}

type EnvVar struct {
	ConfigId int64
	Key      string
	Value    string
}

type Location struct {
	Id          int64
	BaoId       int64
	Subdomain   string
	Destination string
}

type Bao struct {
	Id         int64
	GistId     string `sql:"type:text;"`
	Url        string `sql:"type:text;"`
	Console    string `sql:"type:text;"`
	IsComplete bool
	GitPullUrl string `sql:"type:text;"`
	BaoFileUrl string `sql:"type:text;"`
	Location   Location
	Files      []File
	Config     Config
}

type File struct {
	Id       int64
	BaoId    int64
	Filename string
	Language string
	RawUrl   string
}

func init() {

	var err error
	DB, err = gorm.Open("postgres", "dbname=gitbaotest sslmode=disable")
	if err != nil {
		panic(err)
	}

	DB.DropTableIfExists(&Config{})
	DB.DropTableIfExists(&EnvVar{})
	DB.DropTableIfExists(&Location{})
	DB.DropTableIfExists(&Bao{})
	DB.DropTableIfExists(&File{})

	DB.AutoMigrate(&Location{}, &Bao{}, &File{}, &EnvVar{}, &Config{})
}
