package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gitbao/gitbao/github"
	"github.com/gitbao/gitbao/model"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/{username}/{gist-id}", DownloadHandler).Methods("GET").Host("{subdomain:gist}.{host:.*}")
	r.HandleFunc("/poll/{id}/{line-count}/", PollHandler).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}

var T = template.Must(template.ParseGlob("templates/*"))

func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	T.ExecuteTemplate(w, tmpl+".html", data)
}

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	RenderTemplate(w, "index", nil)
}

func DownloadHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	gistId := vars["gist-id"]
	username := vars["username"]

	path := "/" + username + "/" + gistId

	bao := model.Bao{
		GistId: gistId,
		Console: "Welcome to gitbao!!\n" +
			"Getting ready to wrap up a tasty new bao.\n",
	}

	// go func() {
	err := github.GetGistData(&bao, path, false)
	if err != nil {
		fmt.Printf("%#v", bao)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	stringId := strconv.Itoa(int(bao.Id))
	bao.Location = model.Location{
		Subdomain: bao.GistId + "-" + stringId,
	}

	bao.Console += "Found some files:\n"

	var isGo bool
	for _, value := range bao.Files {
		bao.Console += "    " + value.Filename + "\n"
		if value.Language == "Go" {
			isGo = true
		}
	}

	if isGo != true {
		bao.Console += "Whoops!\n" +
			"gitbao only supports Go programs at the moment.\n" +
			"Quitting...."
		bao.IsComplete = true
	} else {
		// hit up the ziaolong
	}
	// }()

	// Remove this
	bao.IsComplete = true

	query := model.DB.Create(&bao)
	if query.Error != nil {
		fmt.Printf("%#v", bao)
		panic(query.Error)
	}
	RenderTemplate(w, "bao", bao)

}

type pollResponse struct {
	Subdomain  string
	Console    string
	IsComplete bool
}

func PollHandler(w http.ResponseWriter, req *http.Request) {
	time.Sleep(time.Second * 1)

	vars := mux.Vars(req)

	id := vars["id"]
	intid, err := strconv.Atoi(id)
	int64id := int64(intid)

	if err != nil || id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var bao model.Bao
	model.DB.Find(&bao, int64id)

	response := pollResponse{
		IsComplete: bao.IsComplete,
		Console:    bao.Console,
	}

	responseJson, err := json.Marshal(response)
	w.Write(responseJson)
	return
}

func RouterHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("hello"))
}
