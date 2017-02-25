package handler

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/FarmersFriendService/service"
)

func ListFarmers(service service.FarmerService) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		farmers, e := service.ListFarmers()
		if e != "" {
			res.WriteHeader(500)
			return
		}

		farmerDetails, err := json.Marshal(farmers)
		if err != nil {
			res.WriteHeader(500)
			return
		}

		res.Write([]byte(farmerDetails))
	}
}

func AddFarmer(service service.FarmerService) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		
		farmerJson, err := ioutil.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(500)
		}
		
		added, e := service.AddFarmer(farmerJson)
		//todo: can we send back the response body for 500 ?
		if !added {
			fmt.Println("unable to persist farmer: ", string(farmerJson), "error: ", e)
			res.WriteHeader(500)
		}
		res.WriteHeader(200)
	}
}
//
//func GetFarmer(db *sql.DB) http.HandlerFunc {
//	return func(res http.ResponseWriter, req *http.Request) {
//		vars := mux.Vars(req)
//		id := vars["id"]
//		farmerId, _ := strconv.Atoi(id)
//
//		rows, err := db.Query("SELECT * FROM farmers where farmerId = ?", farmerId)
//
//		if err != nil {
//			panic(err)
//		}
//		var farmer service.Farmer
//
//		for rows.Next() {
//			err = rows.Scan(&farmer.Id, &farmer.Name, &farmer.District, &farmer.State, &farmer.PhoneNumber, &farmer.IsDeleted)
//			if err != nil {
//				panic(err)
//			}
//		}
//
//		farmerDetails, err := json.Marshal(farmer)
//		if err != nil {
//			panic(err)
//		}
//
//		res.Write([]byte(farmerDetails))
//	}
//}
//
//func DeleteFarmer(db *sql.DB) http.HandlerFunc {
//	return func(res http.ResponseWriter, req *http.Request) {
//		vars := mux.Vars(req)
//		id := vars["id"]
//		farmerId, _ := strconv.Atoi(id)
//
//		transaction, err := db.Begin()
//		defer func() {
//			switch err {
//			case nil:
//				err = transaction.Commit()
//			default:
//				transaction.Rollback()
//			}
//		}()
//
//		if err != nil {
//			res.WriteHeader(500)
//		}
//		result, err := transaction.Exec("UPDATE farmers SET isDeleted = 1 WHERE farmerId = ?", farmerId)
//
//		if val, err := result.RowsAffected(); val != 1 || err != nil {
//			res.WriteHeader(500)
//		}
//		res.WriteHeader(200)
//	}
//}
//
