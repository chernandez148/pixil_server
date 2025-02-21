package handler

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", Handler) // No need to import handler
	http.ListenAndServe(":8080", nil)
}
