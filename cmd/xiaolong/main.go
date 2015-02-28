package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gitbao/gitbao/model"
	"github.com/gorilla/mux"
)

var server model.Server

func main() {
	if model.Env == "production" {
		myId := os.Getenv("SERVER_ID")
		myIdInt, err := strconv.Atoi(myId)
		if err != nil {
			panic(err)
		}

		query := model.DB.Find(&server, myIdInt)
		if query.Error != nil {
			panic(query.Error)
		}
	} else {
		query := model.DB.First(&server)
		if query.Error != nil {
			panic(query.Error)
		}
	}

	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/ready/{bao-id}", ReadyHandler).Methods("GET")
	r.HandleFunc("/build/{bao-id}", BuildHandler).Methods("GET")
	r.HandleFunc("/logs/{bao-id}", LogHandler).Methods("GET")
	http.Handle("/", Middleware(r))
	fmt.Println("Broadcasting Xiaolong on port 8002")
	http.ListenAndServe(":8002", nil)
}

func Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Host, r.URL)
		h.ServeHTTP(w, r)
	})
}

func ReadyHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	baoIdString := vars["bao-id"]
	baoId, err := strconv.Atoi(baoIdString)
	if baoIdString == "" || err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error parsing bao id"))
		return
	}
	var bao model.Bao
	model.DB.Find(&bao, int64(baoId))

	if bao.IsComplete == true {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Bao is already complete"))
		return
	}

	bao.IsReady = true
	model.DB.Save(&bao)
	return
}

func BuildHandler(w http.ResponseWriter, req *http.Request) {
	// go func() {
	// 	err := builder.StartBuild(&bao, server)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }()

}

func LogHandler(w http.ResponseWriter, req *http.Request) {
	return
}
