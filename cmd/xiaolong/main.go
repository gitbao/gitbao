package main

import (
	"net/http"
	"strconv"

	"github.com/gitbao/gitbao/builder"
	"github.com/gitbao/gitbao/model"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/build/{bao-id}", BuildHandler).Methods("GET")
	r.HandleFunc("/logs/{bao-id}", LogHandler).Methods("GET")
	http.Handle("/", r)
	http.ListenAndServe(":8002", nil)
}

func BuildHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	baoIdString := vars["bao-id"]
	baoId, err := strconv.Atoi(baoIdString)
	if baoIdString == "" || err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error parsing bao id"))
		return
	}
	var bao model.Bao
	model.DB.Find(&bao, int64(baoId))

	bao.Location.Destination = "localhost:8000"

	go func() {
		err := builder.StartBuild(&bao)
		if err != nil {
			panic(err)
		}
	}()
	return
}

func LogHandler(w http.ResponseWriter, req *http.Request) {
	return
}
