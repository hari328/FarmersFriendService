package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/gorilla/mux"
	"net/http"
	"database/sql"
	"github.com/FarmersFriendService/api"
)

func main() {
	router := mux.NewRouter()

	db := initDb("./farmerApp.db")
	listFarmersHandler := api.ListFarmers(db)
	addFarmersHandler := api.AddFarmer(db)
	getFarmersHandler := api.GetFarmer(db)
	deleteFarmerHandler := api.DeleteFarmer(db)


	router.HandleFunc("/farmers", listFarmersHandler).Methods("GET")
	router.HandleFunc("/farmers/{id:[0-9]+}", getFarmersHandler).Methods("GET")
	router.HandleFunc("/farmers", addFarmersHandler).Methods("POST")
	router.HandleFunc("/farmers/{id:[0-9]+}", deleteFarmerHandler).Methods("PATCH")

	http.Handle("/", router)
	http.ListenAndServe(":7000", nil)
}

func initDb(dbName string) *sql.DB{
	db, err := sql.Open("sqlite3", dbName)
	if err != nil { panic(err) }
	if db == nil { panic("db nil") }
	return db
}
