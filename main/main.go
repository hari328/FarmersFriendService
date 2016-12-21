package main

import (
	"database/sql"
	"github.com/golang/go/src/pkg/fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/FarmersFriendService/api"
	"github.com/gorilla/mux"
	"net/http"
)


func ReadItem(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM farmers")
	if err != nil {
		panic(err)
	}
	var fid int
	var name string
	var dist string
	var state string
	var pn string

	for rows.Next() {
		err = rows.Scan(&fid, &name, &dist, &state, &pn)
		if err != nil {
			panic(err)
		}
		fmt.Println(fid)
		fmt.Println(name)
		fmt.Println(dist)
		fmt.Println(state)
		fmt.Println(pn)
	}

	rows.Close()

}
func main() {
	router := mux.NewRouter()

	app := api.NewApi("./farmerApp.db")

	farmersRoute := router.PathPrefix("/farmers")
	farmersRoute.HandlerFunc(app.ListFarmers)

	http.Handle("/", router)
	http.ListenAndServe(":7000", nil)
}

//
//func helloWorld(w http.ResponseWriter, r *http.Request){
//	w.Write([]byte("hello world"))
//}