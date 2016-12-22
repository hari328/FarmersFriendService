package api

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"fmt"
)

func (api *Api) ListFarmers(w http.ResponseWriter, r *http.Request) {

	rows, err := api.Db.Query("SELECT * FROM farmers")
	if err != nil {
		panic(err)
	}

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


	stmt, err := api.Db.Prepare("INSERT INTO farmers(name, district, state, phoneNumber) VALUES (?, ?, ?, ?)")
	if err != nil {
		w.WriteHeader(500)
	}

	fmt.Println(farmer, err, stmt)
	res, err := stmt.Exec(farmer.Name, &farmer.District, farmer.State, farmer.PhoneNumber)
	if err != nil {
		w.WriteHeader(500)
	}

	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)

	_, err = res.LastInsertId()
	if err != nil {
		panic(err)
	}


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
