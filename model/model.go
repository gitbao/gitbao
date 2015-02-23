package model

import (
	"fmt"
	"os"

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

type Server struct {
	Id         int64
	Ip         string
	InstanceId string
	Kind       string
	Dockers    []Docker
}

type Docker struct {
	Id int64
}

func init() {
	Connect()
}
func Connect() {
	var err error

	environment := os.Getenv("GO_ENV")
	if environment == "production" {
		port := "5432"
		host := os.Getenv("GITBAO_DBHOST")
		dbname := os.Getenv("GITBAO_DBNAME")
		username := os.Getenv("GITBAO_DBUSERNAME")
		password := os.Getenv("GITBAO_DBPASSWORD")
		configString := "host=" + host + " port=" + port + " user=" + username + " password=" + password + " sslmode=disable dbname=" + dbname
		fmt.Println(configString)
		DB, err = gorm.Open("postgres", configString)
	} else {
		DB, err = gorm.Open("postgres", "dbname=gitbaotest sslmode=disable")
	}
	if err != nil {
		panic(err)
	}

	tables := []interface{}{
		&Config{},
		&EnvVar{},
		&Location{},
		&Bao{},
		&File{},
		&Server{},
		&Docker{},
	}

	if environment != "production" {
		for _, value := range tables {
			DB.DropTableIfExists(value)
		}
	}

	DB.AutoMigrate(tables...)
}

func Close() {
	err := DB.DB().Close()
	if err != nil {
		panic(err)
	}
}
