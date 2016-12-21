package api

import (
	"net/http"
	"encoding/json"
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
