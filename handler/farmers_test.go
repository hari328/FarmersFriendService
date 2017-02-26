package handler

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/golang/go/src/pkg/strings"

	"github.com/FarmersFriendService/mocks"
	"github.com/FarmersFriendService/util"
	"github.com/FarmersFriendService/model"
	"fmt"
)

const (
	ListFarmersHandler = "ListFarmers"
	AddFarmerHandler = "AddFarmer"
	GerFarmerHandler = "GetFarmer"
	DeleteFarmerHandler = "DeleteFarmer"
)

func TestShouldFetchFarmers(t *testing.T) {
	farmerService := &mocks.MockFarmerService{}
	farmerService.On(ListFarmersHandler).Return(util.GetDummyFarmers(), "")
	
	router := setupRouterForListFarmer(farmerService)
	req, _ := http.NewRequest("GET", "http://localhost/farmers/", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res,req)
	
	result := make([]model.Farmer,0)
	_ = json.Unmarshal(res.Body.Bytes(), &result)
	
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, util.GetDummyFarmers(), result)
}

func TestShouldNotFetchFarmersOnError(t *testing.T) {
	farmerService := &mocks.MockFarmerService{}
	farmerService.On(ListFarmersHandler).Return( nil, "error")
	
	router := setupRouterForListFarmer(farmerService)
	req, _ := http.NewRequest("GET", "http://localhost/farmers/", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res,req)
	
	assert.Equal(t, 500, res.Code)
}

func setupRouterForListFarmer(service *mocks.MockFarmerService) *mux.Router {
	router := mux.NewRouter()
	farmerHandler := ListFarmers(service)
	router.HandleFunc("/farmers/", farmerHandler).Methods("GET")
	return router
}

func TestShouldAddFarmerOnValidInput(t *testing.T) {
	farmerService := &mocks.MockFarmerService{}
	farmerService.On(AddFarmerHandler, []byte("something")).Return(true, "")
	
	router := setupRouterForAddFarmer(farmerService)
	req, _ := http.NewRequest("POST", "http://localhost/farmers/", strings.NewReader("something"))
	res := httptest.NewRecorder()
	router.ServeHTTP(res,req)
	
	assert.Equal(t, 200, res.Code)
}

func TestShouldNotAddFarmerOnError(t *testing.T) {
	farmerService := &mocks.MockFarmerService{}
	farmerService.On(AddFarmerHandler, []byte("something")).Return( false, "db error")
	
	router := setupRouterForAddFarmer(farmerService)
	req, _ := http.NewRequest("POST", "http://localhost/farmers/", strings.NewReader("something"))
	res := httptest.NewRecorder()
	router.ServeHTTP(res,req)
	
	assert.Equal(t, 500, res.Code)
}

func TestShouldNotAddFarmerOnInvalidInput(t *testing.T) {
	farmerService := &mocks.MockFarmerService{}
	farmerService.On(AddFarmerHandler, []byte("something")).Return(false, "json couldn't be decoded")
	
	router := setupRouterForAddFarmer(farmerService)
	
	req, _ := http.NewRequest("POST", "http://localhost/farmers/", strings.NewReader("something"))
	res := httptest.NewRecorder()
	router.ServeHTTP(res,req)
	
	assert.Equal(t, 500, res.Code)
}

func setupRouterForAddFarmer(service *mocks.MockFarmerService) *mux.Router {
	router := mux.NewRouter()
	handler := AddFarmer(service)
	router.HandleFunc("/farmers/", handler).Methods("POST")
	return router
}

func TestShouldGetParticularFarmerDetails(t *testing.T) {
	farmerService := &mocks.MockFarmerService{}
	farmerService.On(GerFarmerHandler, 1).Return(util.DummyFarmerOne(), "")
	
	router := setupRouterForGetFarmer(farmerService)
	
	req, _ := http.NewRequest("GET", "http://localhost/farmers/1", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	
	result := model.Farmer{}
	_ = json.Unmarshal(res.Body.Bytes(), &result)
	
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, util.DummyFarmerOne(), result)
}
func TestShouldNotGetParticularFarmerDetailsIfNotPresent(t *testing.T) {
	farmerService := &mocks.MockFarmerService{}
	farmerService.On(GerFarmerHandler, 1).Return(model.Farmer{}, "")
	
	router := setupRouterForGetFarmer(farmerService)
	
	req, _ := http.NewRequest("GET", "http://localhost/farmers/1", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res,req)
	
	result := model.Farmer{}
	_ = json.Unmarshal(res.Body.Bytes(), &result)
	
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, model.Farmer{}, result)
}

func TestShouldNotGetParticularFarmerDetailsOnError(t *testing.T) {
	farmerService := &mocks.MockFarmerService{}
	farmerService.On(GerFarmerHandler, 1).Return(model.Farmer{}, "db Error")
	
	router := setupRouterForGetFarmer(farmerService)
	
	req, _ := http.NewRequest("GET", "http://localhost/farmers/1", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	
	assert.Equal(t, 500, res.Code)
}

func setupRouterForGetFarmer(farmerService *mocks.MockFarmerService) *mux.Router {
	router := mux.NewRouter()
	handler := GetFarmer(farmerService)
	router.HandleFunc("/farmers/{id:[0-9]+}", handler).Methods("GET")
	return router
}


func TestShouldDeleteParticularFarmer(t *testing.T) {
	farmerService := &mocks.MockFarmerService{}
	farmerService.On(DeleteFarmerHandler, 1).Return(nil)
	
	router := setupRouterForDeleteFarmer(farmerService)
	
	req, _ := http.NewRequest("PATCH", "http://localhost/farmers/1", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res,req)
	
	assert.Equal(t, 200, res.Code)
}

func TestShouldNotDeleteFarmerOnError(t *testing.T) {
	farmerService := &mocks.MockFarmerService{}
	farmerService.On(DeleteFarmerHandler, 1).Return(fmt.Errorf("unable to find record"))
	
	router := setupRouterForDeleteFarmer(farmerService)
	
	req, _ := http.NewRequest("PATCH", "http://localhost/farmers/1", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	
	result := model.Farmer{}
	_ = json.Unmarshal(res.Body.Bytes(), &result)
	
	assert.Equal(t, 500, res.Code)
}

func setupRouterForDeleteFarmer(mockFarmerService *mocks.MockFarmerService) *mux.Router {
	router := mux.NewRouter()
	handler := DeleteFarmer(mockFarmerService)
	router.HandleFunc("/farmers/{id:[0-9]+}", handler).Methods("PATCH")
	return router
}
