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
	ListFarmersHandler  = "ListFarmers"
	AddFarmerHandler    = "AddFarmer"
	GerFarmerHandler    = "GetFarmer"
	DeleteFarmerHandler = "DeleteFarmer"
)

func TestListFarmersReturnsListOfFarmersWhenNoError(t *testing.T) {
	farmerService := mockFarmerRepoForListFarmers(util.GetDummyFarmers(), nil)
	res := makeListFarmersCall(farmerService)
	result := make([]model.Farmer, 0)
	_ = json.Unmarshal(res.Body.Bytes(), &result)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, util.GetDummyFarmers(), result)
}

func TestListFarmersReturnsErrorOnServiceError(t *testing.T) {
	farmerService := mockFarmerRepoForListFarmers(nil, fmt.Errorf("error"))
	res := makeListFarmersCall(farmerService)
	assert.Equal(t, 500, res.Code)
}

func mockFarmerRepoForListFarmers(returnValue []model.Farmer, returnError error) *mocks.MockFarmerService {
	farmerService := &mocks.MockFarmerService{}
	farmerService.On(ListFarmersHandler).Return(returnValue, returnError)
	return farmerService
}

func makeListFarmersCall(farmerService *mocks.MockFarmerService) *httptest.ResponseRecorder {
	router := setupRouterForListFarmer(farmerService)
	req, _ := http.NewRequest("GET", "http://localhost/farmers/", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}

func setupRouterForListFarmer(service *mocks.MockFarmerService) *mux.Router {
	router := mux.NewRouter()
	farmerHandler := ListFarmers(service)
	router.HandleFunc("/farmers/", farmerHandler).Methods("GET")
	return router
}

func TestAddFarmerReturnSuccessOnValidInput(t *testing.T) {
	requestBody := "something"
	farmerService := mockFarmerServiceForAddFarmer(requestBody, nil)
	res := makeAddFarmerCall(farmerService, requestBody)
	assert.Equal(t, 200, res.Code)
}

func TestAddFarmerReturnErrorWhenInputIsEmpty(t *testing.T) {
	res := makeAddFarmerCall(&mocks.MockFarmerService{}, "")
	assert.Equal(t, 500, res.Code)
}

func TestShouldNotAddFarmerOnError(t *testing.T) {
	requestBody := "something"
	farmerService := mockFarmerServiceForAddFarmer(requestBody, fmt.Errorf("db error"))
	res := makeAddFarmerCall(farmerService, requestBody)
	assert.Equal(t, 500, res.Code)
}

func mockFarmerServiceForAddFarmer(requestBody string, returnError error) *mocks.MockFarmerService {
	farmerService := &mocks.MockFarmerService{}
	farmerService.On(AddFarmerHandler, []byte(requestBody)).Return(returnError)
	return farmerService
}

func makeAddFarmerCall(service *mocks.MockFarmerService, requestBody string) *httptest.ResponseRecorder {
	router := setupRouterForAddFarmer(service)
	req, _ := http.NewRequest("POST", "http://localhost/farmers/", strings.NewReader(requestBody))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}

func setupRouterForAddFarmer(service *mocks.MockFarmerService) *mux.Router {
	router := mux.NewRouter()
	handler := AddFarmer(service)
	router.HandleFunc("/farmers/", handler).Methods("POST")
	return router
}

func TestGetFarmerReturnsFarmerOnValidInput(t *testing.T) {
	farmerService := mockFarmerServiceForGetFarmer(1, util.DummyFarmerOne(), nil)
	res := makeGetFarmerCall(farmerService, "1")
	
	result := model.Farmer{}
	_ = json.Unmarshal(res.Body.Bytes(), &result)
	
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, util.DummyFarmerOne(), result)
}

func TestGetFarmerReturnsErrorOnServiceError(t *testing.T) {
	farmerService := mockFarmerServiceForGetFarmer(1, model.Farmer{}, fmt.Errorf("error"))
	res := makeGetFarmerCall(farmerService, "1")

	result := model.Farmer{}
	_ = json.Unmarshal(res.Body.Bytes(), &result)
	
	assert.Equal(t, 500, res.Code)
}

func mockFarmerServiceForGetFarmer(id int, resultValue model.Farmer, err error) *mocks.MockFarmerService {
	farmerService := &mocks.MockFarmerService{}
	farmerService.On(GerFarmerHandler, id).Return(resultValue, err)
	return farmerService
}

func setupRouterForGetFarmer(farmerService *mocks.MockFarmerService) *mux.Router {
	router := mux.NewRouter()
	handler := GetFarmer(farmerService)
	router.HandleFunc("/farmers/{id:[0-9]+}", handler).Methods("GET")
	return router
}

func makeGetFarmerCall(service *mocks.MockFarmerService, id string) *httptest.ResponseRecorder {
	router := setupRouterForGetFarmer(service)
	getFarmerUrl := fmt.Sprintf("http://localhost/farmers/%s", id)
	req, _ := http.NewRequest("GET", getFarmerUrl , nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}

func TestDeleteFarmerReturnSuccessOnValidInput(t *testing.T) {
	farmerService := mockFarmerServiceForDeleteFarmer(1, nil)
	res := makeDeleteFarmerCall(farmerService, "1")
	assert.Equal(t, 200, res.Code)
}

func TestDeleteFarmerReturnErrorOnServiceError(t *testing.T) {
	farmerService := mockFarmerServiceForDeleteFarmer(1,fmt.Errorf("error"))
	
	res := makeDeleteFarmerCall(farmerService, "1")

	result := model.Farmer{}
	_ = json.Unmarshal(res.Body.Bytes(), &result)

	assert.Equal(t, 500, res.Code)
}

func mockFarmerServiceForDeleteFarmer(id int, err error) *mocks.MockFarmerService {
	farmerService := &mocks.MockFarmerService{}
	farmerService.On(DeleteFarmerHandler, id).Return(err)
	return farmerService
}

func setupRouterForDeleteFarmer(farmerService *mocks.MockFarmerService) *mux.Router {
	router := mux.NewRouter()
	handler := DeleteFarmer(farmerService)
	router.HandleFunc("/farmers/{id:[0-9]+}", handler).Methods("PATCH")
	return router
}

func makeDeleteFarmerCall(service *mocks.MockFarmerService, id string) *httptest.ResponseRecorder {
	router := setupRouterForDeleteFarmer(service)
	deleteFarmerUrl := fmt.Sprintf("http://localhost/farmers/%s", id)
	
	req, _ := http.NewRequest("PATCH", deleteFarmerUrl, nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}
