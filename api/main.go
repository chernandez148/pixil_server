package main

import (
	"net/http"
	"pixi/api/handler"
)

func main() {
	http.HandleFunc("/", handler.Handler)
	http.ListenAndServe(":8080", nil)
}
