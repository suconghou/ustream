package main

import (
	"net/http"
	"now/api"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	serve(port)
}

func serve(port string) error {
	http.HandleFunc("/", api.Handler)
	return http.ListenAndServe(":"+port, nil)
}
