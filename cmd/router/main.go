package main

import (
	"net/http"

	"github.com/gitbao/gitbao/router"
)

func main() {
	http.ListenAndServe(":8001", &router.Router{})
}
