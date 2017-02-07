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

	registerFarmerRoutes(db, router)

	http.Handle("/", router)
	http.ListenAndServe(":9000", nil)
}

func initDb(dbName string) *sql.DB{
	db, err := sql.Open("sqlite3", dbName)
	if err != nil { panic(err) }
	if db == nil { panic("db nil") }
	return db
}

func registerFarmerRoutes(db *sql.DB, rootRouter *mux.Router ) {
	farmersRouter := rootRouter.PathPrefix("/farmers").Subrouter()

	listFarmersHandler := api.ListFarmers(db)
	addFarmersHandler := api.AddFarmer(db)
	getFarmersHandler := api.GetFarmer(db)
	deleteFarmerHandler := api.DeleteFarmer(db)

	farmersRouter.HandleFunc("/", listFarmersHandler).Methods("GET")
	farmersRouter.HandleFunc("/{id:[0-9]+}", getFarmersHandler).Methods("GET")
	farmersRouter.HandleFunc("/", addFarmersHandler).Methods("POST")
	farmersRouter.HandleFunc("/{id:[0-9]+}", deleteFarmerHandler).Methods("PATCH")
}