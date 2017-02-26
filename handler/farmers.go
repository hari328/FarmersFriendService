package handler

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/FarmersFriendService/service"
	"github.com/gorilla/mux"
	"strconv"
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

func GetFarmer(service service.FarmerService) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]
		farmerId, _ := strconv.Atoi(id)

		farmer, er := service.GetFarmer(farmerId)
		
		if er != "" {
			fmt.Println("unable to get farmer for id :",id , er)
			res.WriteHeader(500)
		}
		
		farmerDetails, err := json.Marshal(farmer)
		if err != nil {
			fmt.Println("failed to marshal json", err.Error())
			res.WriteHeader(500)
		}
		
		res.WriteHeader(200)
		res.Write([]byte(farmerDetails))
	}
}

func DeleteFarmer(service service.FarmerService) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := vars["id"]
		farmerId, err := strconv.Atoi(id)
		if err != nil {
			res.WriteHeader(500)
		}
		err = service.DeleteFarmer(farmerId)
		
		if err != nil {
			fmt.Println("unable to delete farmer",err)
			res.WriteHeader(500)
		}
		res.WriteHeader(200)
	}
}

//transaction, err := db.Begin()
//defer func() {
//	switch err {
//	case nil:
//		err = transaction.Commit()
//	default:
//		transaction.Rollback()
//	}
//}()