package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/hello", helloWorld)
	http.Handle("/", router)
	http.ListenAndServe(":7000", nil)
}

func helloWorld(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Hello World"))
}