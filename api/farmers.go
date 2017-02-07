package api

import (
	"net/http"
	"encoding/json"
	"github.com/FarmersFriendService/model"
	"database/sql"
	"io/ioutil"
	"github.com/gorilla/mux"
	"strconv"
)

func ListFarmers(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		rows, err := db.Query("SELECT * FROM farmers WHERE isDeleted = 0")
		if err != nil {
			panic(err)
		}

		farmers := &Farmers{List: make([]model.Farmer, 0)}

		var farmer model.Farmer
		for rows.Next() {
			err = rows.Scan(&farmer.Id, &farmer.Name, &farmer.District, &farmer.State, &farmer.PhoneNumber, &farmer.IsDeleted)
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

		res.Write([]byte(farmerDetails))
	}
}

func AddFarmer(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {

		farmerJson, err := ioutil.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(500)
		}

		var farmer model.Farmer
		if err := json.Unmarshal(farmerJson, &farmer); err != nil {
			res.WriteHeader(500)
		}

		transaction, err := db.Begin()
		defer func() {
			switch err {
			case nil:
				err = transaction.Commit()
			default:
				transaction.Rollback()
			}
		}()

		if err != nil {
			res.WriteHeader(500)
		}

		isDeleted := 0

		result, err := transaction.Exec("INSERT INTO farmers(name, district, state, phoneNumber, isDeleted) VALUES (?, ?, ?, ?, ?)", farmer.Name, farmer.District, farmer.State, farmer.PhoneNumber, isDeleted)

		if err != nil {
			res.WriteHeader(500)
		}

		if val, err := result.RowsAffected(); val != 1 || err != nil {
			res.WriteHeader(500)
		}

		res.WriteHeader(200)
	}
}

func GetFarmer(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]
		farmerId, _ := strconv.Atoi(id)

		rows, err := db.Query("SELECT * FROM farmers where farmerId = ?", farmerId)

		if err != nil {
			panic(err)
		}
		var farmer model.Farmer

		for rows.Next() {
			err = rows.Scan(&farmer.Id, &farmer.Name, &farmer.District, &farmer.State, &farmer.PhoneNumber, &farmer.IsDeleted)
			if err != nil {
				panic(err)
			}
		}

		farmerDetails, err := json.Marshal(farmer)
		if err != nil {
			panic(err)
		}

		res.Write([]byte(farmerDetails))
	}
}

func DeleteFarmer(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]
		farmerId, _ := strconv.Atoi(id)

		transaction, err := db.Begin()
		defer func() {
			switch err {
			case nil:
				err = transaction.Commit()
			default:
				transaction.Rollback()
			}
		}()

		if err != nil {
			res.WriteHeader(500)
		}
		result, err := transaction.Exec("UPDATE farmers SET isDeleted = 1 WHERE farmerId = ?", farmerId)

		if val, err := result.RowsAffected(); val != 1 || err != nil {
			res.WriteHeader(500)
		}
		res.WriteHeader(200)
	}
}

type Farmers struct {
	List []model.Farmer 			`json:"farmers"`
}
