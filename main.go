package main

import (
	"net/http"

	"github.com/gitbao/gitbao/github"
	"github.com/gitbao/gitbao/model"
	"github.com/sqs/mux"
)

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/{username}/{gist-id}", DownloadHandler).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8000", nil)
}

func DownloadHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	gistId := vars["gist-id"]
	username := vars["username"]

	path := "/" + username + "/" + gistId

	bao, err := github.GetGistData(gistId, path, false)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(bao.GitPullUrl))
	model.DB.Create(&bao)
}
