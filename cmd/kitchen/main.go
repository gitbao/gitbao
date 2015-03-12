package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gitbao/gitbao/github"
	"github.com/gitbao/gitbao/model"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var T *template.Template
var goPath string

func init() {
	goPath = os.Getenv("GOPATH")
	T = template.Must(template.ParseGlob(goPath + "/src/github.com/gitbao/gitbao/cmd/kitchen/templates/*"))
}
func main() {

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/", IndexHandler).Methods("GET")
	r.HandleFunc("/{username}/{gist-id}", CreateHandler).Methods("GET").Host("{subdomain:gist}.{host:.*}")
	r.HandleFunc("/bao/{id-base36}/", BaoHandler).Methods("GET")
	r.HandleFunc("/poll/{id}/", PollHandler).Methods("GET")
	r.HandleFunc("/deploy/{id}/", DeployHandler).Methods("POST")
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

func CreateHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	gistId := vars["gist-id"]
	username := vars["username"]
	log.Printf("New bao: %s %s\n", gistId, username)

	if gistId == "" || username == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404"))
		return
	}

	bao := model.Bao{
		GistId: gistId,
	}

	query := model.DB.Create(&bao)
	if query.Error != nil {
		log.Println("Error:", query.Error)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(query.Error.Error()))
		return
	}

	host := req.Host
	host_parts := strings.Split(host, ".")

	base36Id := strconv.FormatInt(bao.Id, 36)

	http.Redirect(
		w, req,
		fmt.Sprintf(
			"http://%s/bao/%s/",
			strings.Join(host_parts[1:], "."),
			base36Id,
		),
		302,
	)
}

func BaoHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	base36Id := vars["id-base36"]
	fmt.Println(base36Id)
	baoId, err := strconv.ParseInt(base36Id, 36, 64)

	var bao model.Bao
	query := model.DB.Find(&bao, baoId)
	if query.Error == gorm.RecordNotFound {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404"))
		return
	} else if query.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if bao.GitPullUrl == "" {

		// go func() {
		err = github.GetGistData(&bao)
		if err != nil {
			fmt.Printf("%#v", bao)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		bao.Console = "Welcome to gitbao!!\n" +
			"Getting ready to wrap up a tasty new bao.\n"

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
			bao.Console += "Nice, looks like we can deploy your application\n" +
				"Modify your config file if needed and hit deploy!\n"
		}

		if bao.IsComplete != true {
			go func() {
				var server model.Server
				query = model.DB.Where("kind = ?", "xiaolong").Find(&server)
				if query.Error != nil {
					bao.Console += "Uh oh, we've experienced an error. Please try again.\n"
					fmt.Println(query.Error)
					model.DB.Save(&bao)
					return
				}
				bao.ServerId = server.Id
				model.DB.Save(&bao)
				getUrl := fmt.Sprintf("http://%s:8002/ready/%d", server.Ip, bao.Id)
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
		}
	}

	query = model.DB.Save(&bao)
	if query.Error != nil {
		fmt.Printf("%#v", bao)
		panic(query.Error)
	}
	RenderTemplate(w, "bao", bao)

}

func DeployHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	id := vars["id"]
	intid, err := strconv.Atoi(id)
	int64id := int64(intid)

	if err != nil || id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var bao model.Bao
	query := model.DB.Find(&bao, int64id)
	if query.Error == gorm.RecordNotFound {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404"))
		return
	} else if query.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	var server model.Server
	query = model.DB.Find(&server, bao.ServerId)
	if query.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
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

}

type pollResponse struct {
	Subdomain  string
	Console    string
	IsReady    bool
	IsComplete bool
	Url        string
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

	var location model.Location
	model.DB.Find(&location, bao.Id)

	var url string
	if location.Subdomain != "" {
		url = location.Subdomain + ".gitbao.com"
	}
	response := pollResponse{
		IsComplete: bao.IsComplete,
		Console:    bao.Console,
		IsReady:    bao.IsReady,
		Url:        url,
	}

	responseJson, err := json.Marshal(response)
	w.Write(responseJson)
	return
}

func RouterHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("hello"))
}
