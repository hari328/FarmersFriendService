package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/FarmersFriendService/handler"
	"github.com/gocraft/dbr"
	"github.com/FarmersFriendService/repository"
)

func main() {
	router := mux.NewRouter()
	db := initDb("./farmerApp.db")
	farmerService := repository.New(db)
	
	registerFarmerRoutes(farmerService, router)

	http.Handle("/", router)
	http.ListenAndServe(":9000", nil)
}

func registerFarmerRoutes(repository repository.FarmerRepository, rootRouter *mux.Router ) {
	listFarmersHandler := handler.ListFarmers(repository)
	addFarmersHandler := handler.AddFarmer(repository)
	getFarmersHandler := handler.GetFarmer(repository)
	deleteFarmerHandler := handler.DeleteFarmer(repository)
	
	farmersRouter := rootRouter.PathPrefix("/farmers").Subrouter()
	
	farmersRouter.HandleFunc("/", listFarmersHandler).Methods("GET")
	farmersRouter.HandleFunc("/{id:[0-9]+}", getFarmersHandler).Methods("GET")
	farmersRouter.HandleFunc("/", addFarmersHandler).Methods("POST")
	farmersRouter.HandleFunc("/{id:[0-9]+}", deleteFarmerHandler).Methods("PATCH")
}

func initDb(dbName string) *dbr.Connection{
	db, err := dbr.Open("sqlite3", dbName, nil)
	if err != nil { panic(err) }
	if db == nil { panic("db nil") }
	return db
}
