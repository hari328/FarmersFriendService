package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/FarmersFriendService/api"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	app := api.NewApi("./farmerApp.db")

	router.HandleFunc("/farmers", app.ListFarmers)
	router.HandleFunc("/farmers/{id:[0-9]+}",app.GetFarmer)

	router.HandleFunc("/problems", app.ListProblems)

	farmersRoutePost := router.PathPrefix("/farmers").Methods("POST")
	farmersRoutePost.HandlerFunc(app.AddFarmer)

	http.Handle("/", router)
	http.ListenAndServe(":7000", nil)
}
