package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/gitbao/gitbao/github"
	"github.com/gitbao/gitbao/model"
	"github.com/gorilla/mux"
)

var T *template.Template

func main() {

	goPath := os.Getenv("GOPATH")
	T = template.Must(template.ParseGlob(goPath + "src/github.com/gitbao/gitbao/cmd/kitchen/templates/*"))

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/{username}/{gist-id}", DownloadHandler).Methods("GET").Host("{subdomain:gist}.{host:.*}")
	r.HandleFunc("/poll/{id}/{line-count}/", PollHandler).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(goPath + "src/github.com/gitbao/gitbao/cmd/kitchen/public/")))
	http.Handle("/", Middleware(r))
	fmt.Println("Broadcasting Kitchen on port 8000")
	http.ListenAndServe(":8000", nil)
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Host, r.URL)
		h.ServeHTTP(w, r)
	})
}

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

	log.Println("New bao", gistId, username)

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
		if value.Filename == "Baofile" || value.Filename == "baofile" {
			bao.BaoFileUrl = value.RawUrl
		}
	}

	if isGo != true {
		bao.Console += "Whoops!\n" +
			"gitbao only supports Go programs at the moment.\n" +
			"Quitting...."
		bao.IsComplete = true
	} else {
	}
	// }()

	query := model.DB.Create(&bao)
	if query.Error != nil {
		fmt.Printf("%#v", bao)
		panic(query.Error)
	}
	go func() {
		var server model.Server
		query := model.DB.Where("kind = ?", "xiaolong").Find(&server)
		if query.Error != nil {
			bao.Console += "Uh oh, we've experienced an error. Please try again.\n"
			fmt.Println(query.Error)
			model.DB.Save(&bao)
			return
		}
		getUrl := fmt.Sprintf("http://%s:8002/build/%d", server.Ip, bao.Id)
		log.Println(getUrl)
		resp, err := http.Get(getUrl)
		log.Printf("%#v", resp)
		if err != nil || resp.StatusCode != 200 {
			bao.Console += "Uh oh, we've experienced an error. Please try again.\n"
			bao.IsComplete = true
			model.DB.Save(&bao)
			return
		}
	}()
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
