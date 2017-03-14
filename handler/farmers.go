package handler

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/FarmersFriendService/repository"
	"github.com/gorilla/mux"
	"strconv"
)

func ListFarmers(repo repository.FarmerRepository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		farmers, e := repo.ListFarmers()
		fmt.Println(e)
		if e != nil {
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

func AddFarmer(repo repository.FarmerRepository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		farmerJson, err := ioutil.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(500)
			return
		}
		
		if len(farmerJson) == 0 {
			res.WriteHeader(500)
			return
		}
		
		err = repo.AddFarmer(farmerJson)
		if err != nil {
			fmt.Println("unable to persist farmer: ", string(farmerJson), "error: ", err)
			res.WriteHeader(500)
		}
		res.WriteHeader(200)
	}
}

func GetFarmer(repo repository.FarmerRepository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		farmerId := getFarmerId(req)
		farmer, er := repo.GetFarmer(farmerId)
		
		if er != "" {
			fmt.Println("unable to get farmer for id :",farmerId , er)
			res.WriteHeader(500)
			return
		}
		
		farmerDetails, err := json.Marshal(farmer)
		if err != nil {
			fmt.Println("failed to marshal json", err.Error())
			res.WriteHeader(500)
			return
		}
		
		res.WriteHeader(200)
		res.Write([]byte(farmerDetails))
	}
}

func DeleteFarmer(repo repository.FarmerRepository) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		farmerId:= getFarmerId(req)
		err := repo.DeleteFarmer(farmerId)
		
		if err != nil {
			fmt.Println("unable to delete farmer", err)
			res.WriteHeader(500)
			return
		}
		
		res.WriteHeader(200)
	}
}

func getFarmerId(req *http.Request) int {
	vars := mux.Vars(req)
	id := vars["id"]
	farmerId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("unable to get farmer id from request", err)
	}
	return farmerId
}
