package api

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"github.com/gorilla/mux"
	"strconv"
	"fmt"
)

func (api *Api) ListFarmers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here")

	rows, err := api.Db.Query("SELECT * FROM farmers")
	if err != nil {
		panic(err)
	}

	fmt.Println("here")
	farmers := &Farmers{List: make([]Farmer, 0)}

	var farmer Farmer
	for rows.Next() {
		err = rows.Scan(&farmer.Id, &farmer.Name, &farmer.District, &farmer.State, &farmer.PhoneNumber)
		if err != nil {
			panic(err)
		}

		farmers.List = append(farmers.List, farmer)
	}

	rows.Close()

	farmerDetails,err := json.Marshal(farmers)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(farmerDetails))
}

func (api *Api) AddFarmer(w http.ResponseWriter, r *http.Request){

	farmerJson, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
	}

	var farmer Farmer
	if err := json.Unmarshal(farmerJson, &farmer); err != nil {
		w.WriteHeader(500)
	}

	transaction, err := api.Db.Begin()
	defer func() {
		switch err {
		case nil:
			err = transaction.Commit()
		default:
			transaction.Rollback()
		}
	}()

	if err != nil {
		w.WriteHeader(500)
	}

	result, err := transaction.Exec("INSERT INTO farmers(name, district, state, phoneNumber) VALUES (?, ?, ?, ?)", farmer.Name, farmer.District, farmer.State, farmer.PhoneNumber)

	if err != nil {
		w.WriteHeader(500)
	}

	if val, err := result.RowsAffected(); val != 1 || err != nil {
		w.WriteHeader(500)
	}

	w.WriteHeader(200)
}

func (api *Api) GetFarmer(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id := vars["id"]
	farmerId, _ := strconv.Atoi(id)

	rows, err := api.Db.Query("SELECT * FROM farmers where farmerId = ?", farmerId)

	if err != nil {
		panic(err)
	}
	var farmer Farmer

	for rows.Next() {
		err = rows.Scan(&farmer.Id, &farmer.Name, &farmer.District, &farmer.State, &farmer.PhoneNumber)
		if err != nil {
			panic(err)
		}
	}

	farmerDetails,err := json.Marshal(farmer)
	if err != nil {
		panic(err)
	}

	w.Write([]byte(farmerDetails))
}

type Farmers struct {
	List []Farmer 			`json:"farmers"`
}

type Farmer struct {
	Id int								`json:"farmerId"`
	Name string						`json:"name"`
	District string				`json:"district"`
	State string					`json:"state"`
	PhoneNumber int64			`json:"phoneNumber"`
}
