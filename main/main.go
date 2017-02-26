package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/gorilla/mux"
	"net/http"
	"database/sql"
	"github.com/FarmersFriendService/service"
	"github.com/FarmersFriendService/handler"
)

func main() {
	router := mux.NewRouter()
	db := initDb("./farmerApp.db")
	farmerService := service.New(db)
	
	registerFarmerRoutes(farmerService, router)

	http.Handle("/", router)
	http.ListenAndServe(":9000", nil)
}

func registerFarmerRoutes(farmerServicer service.FarmerService, rootRouter *mux.Router ) {
	
	farmersRouter := rootRouter.PathPrefix("/farmers").Subrouter()
	
	listFarmersHandler := handler.ListFarmers(farmerServicer)
	addFarmersHandler := handler.AddFarmer(farmerServicer)
	getFarmersHandler := handler.GetFarmer(farmerServicer)
	deleteFarmerHandler := handler.DeleteFarmer(farmerServicer)
	
	farmersRouter.HandleFunc("/", listFarmersHandler).Methods("GET")
	farmersRouter.HandleFunc("/{id:[0-9]+}", getFarmersHandler).Methods("GET")
	farmersRouter.HandleFunc("/", addFarmersHandler).Methods("POST")
	farmersRouter.HandleFunc("/{id:[0-9]+}", deleteFarmerHandler).Methods("PATCH")
}

func initDb(dbName string) *sql.DB{
	db, err := sql.Open("sqlite3", dbName)
	if err != nil { panic(err) }
	if db == nil { panic("db nil") }
	return db
}
