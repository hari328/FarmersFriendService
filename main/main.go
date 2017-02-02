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
	testHandler := test()

	testRouter := router.Host("test").Subrouter()

	testRouter.HandleFunc("/test/", testHandler).Methods("GET")
	router.HandleFunc("/farmers", listFarmersHandler).Methods("GET")
	//router.HandleFunc("/farmers/{id:[0-9]+}",app.GetFarmer)
	//
	//router.HandleFunc("/problems", app.ListProblems)
	//
	//farmersRoutePost := router.PathPrefix("/farmers").Methods("POST")
	//farmersRoutePost.HandlerFunc(app.AddFarmer)


	http.Handle("/", router)
	http.ListenAndServe(":7000", nil)
}

func initDb(dbName string) *sql.DB{
	db, err := sql.Open("sqlite3", dbName)
	if err != nil { panic(err) }
	if db == nil { panic("db nil") }
	return db
}

func test() http.HandlerFunc {
	return func (res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("hey there"))
	}
}