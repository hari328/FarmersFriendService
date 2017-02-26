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
)

const (
	ListFarmersHandler = "ListFarmers"
	AddFarmerHandler = "AddFarmer"
	GerFarmerHandler = "GetFarmer"
)

func TestShouldFetchFarmers(t *testing.T) {
	farmerService := &mocks.MockFarmerService{}
	farmerService.On(ListFarmersHandler).Return(util.GetDummyFarmers(), "")
	
	router := setupRouterForGetFarmer(farmerService)
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
	
	router := setupRouterForGetFarmer(farmerService)
	req, _ := http.NewRequest("GET", "http://localhost/farmers/", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res,req)
	
	assert.Equal(t, 500, res.Code)
}

func setupRouterForGetFarmer(service *mocks.MockFarmerService) *mux.Router {
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
	mockFarmerService := &mocks.MockFarmerService{}
	mockFarmerService.On(GerFarmerHandler, 1).Return(util.DummyFarmerOne(), "")

	router := mux.NewRouter()
	handler := GetFarmer(mockFarmerService)
	router.HandleFunc("/farmers/{id:[0-9]+}", handler).Methods("GET")

	req, _ := http.NewRequest("GET", "http://localhost/farmers/1", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res,req)
	
	result := model.Farmer{}
	_ = json.Unmarshal(res.Body.Bytes(), &result)
	
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, util.DummyFarmerOne(), result)
}

func TestShouldNotGetParticularFarmerDetailsIfNotPresent(t *testing.T) {
	mockFarmerService := &mocks.MockFarmerService{}
	mockFarmerService.On(GerFarmerHandler, 1).Return(model.Farmer{}, "")
	
	router := mux.NewRouter()
	handler := GetFarmer(mockFarmerService)
	router.HandleFunc("/farmers/{id:[0-9]+}", handler).Methods("GET")
	
	req, _ := http.NewRequest("GET", "http://localhost/farmers/1", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res,req)
	
	result := model.Farmer{}
	_ = json.Unmarshal(res.Body.Bytes(), &result)
	
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, model.Farmer{}, result)
}

func TestShouldNotGetParticularFarmerDetailsOnError(t *testing.T) {
	mockFarmerService := &mocks.MockFarmerService{}
	mockFarmerService.On(GerFarmerHandler, 1).Return(model.Farmer{}, "db Error")
	
	router := mux.NewRouter()
	handler := GetFarmer(mockFarmerService)
	router.HandleFunc("/farmers/{id:[0-9]+}", handler).Methods("GET")
	
	req, _ := http.NewRequest("GET", "http://localhost/farmers/1", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res,req)
	
	assert.Equal(t, 500, res.Code)
}