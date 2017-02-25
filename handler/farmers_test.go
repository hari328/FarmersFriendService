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


func TestShouldFetchFarmers(t *testing.T) {
	farmerModelMock := &mocks.MockFarmerService{}
	farmerModelMock.On("ListFarmers").Return(util.GetDummyFarmers(), "")
	
	router := mux.NewRouter()
	farmerHandler := ListFarmers(farmerModelMock)
	router.HandleFunc("/farmers/", farmerHandler).Methods("GET")
	
	req, _ := http.NewRequest("GET", "http://localhost/farmers/", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res,req)
	
	result := make([]model.Farmer,0)
	_ = json.Unmarshal(res.Body.Bytes(), &result)
	
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, util.GetDummyFarmers(), result)
}

func TestShouldNotFetchFarmersOnError(t *testing.T) {
	farmerModelMock := &mocks.MockFarmerService{}
	farmerModelMock.On("ListFarmers").Return( nil, "error")
	
	router := mux.NewRouter()
	farmerHandler := ListFarmers(farmerModelMock)
	router.HandleFunc("/farmers/", farmerHandler).Methods("GET")
	
	req, _ := http.NewRequest("GET", "http://localhost/farmers/", nil)
	res := httptest.NewRecorder()
	router.ServeHTTP(res,req)
	
	assert.Equal(t, 500, res.Code)
}

func TestShouldAddFarmerOnValidInput(t *testing.T) {
	farmerModelMock := &mocks.MockFarmerService{}
	farmerModelMock.On("AddFarmer", []byte("something")).Return(true, "")
	
	router := mux.NewRouter()
	handler := AddFarmer(farmerModelMock)
	router.HandleFunc("/farmers/", handler).Methods("POST")
	
	req, _ := http.NewRequest("POST", "http://localhost/farmers/", strings.NewReader("something"))
	res := httptest.NewRecorder()
	router.ServeHTTP(res,req)
	
	if res.Code != 200 {
		t.Fatalf("expected status code to be 200, but got: %d", res.Code)
	}
}

func TestShouldNotAddFarmerOnInvalidInput(t *testing.T) {
	farmerModelMock := &mocks.MockFarmerService{}
	farmerModelMock.On("AddFarmer", []byte("something")).Return(false, "json couldn't be decoded")
	
	router := mux.NewRouter()
	handler := AddFarmer(farmerModelMock)
	router.HandleFunc("/farmers/", handler).Methods("POST")
	
	req, _ := http.NewRequest("POST", "http://localhost/farmers/", strings.NewReader("something"))
	res := httptest.NewRecorder()
	router.ServeHTTP(res,req)
	
	assert.Equal(t, 500, res.Code)
}

//func TestShouldGetParticularFarmerDetails(t *testing.T) {
//	router := mux.NewRouter()
//	db, mock, err := sqlmock.New()
//
//	router.HandleFunc("/farmers/{id:[0-9]+}", GetFarmer(db)).Methods("GET")
//
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//
//	defer db.Close()
//	farmerData := dummyFarmerOne()
//
//	rows := sqlmock.NewRows([]string{"id", "name", "district", "state", "phoneNumber", "isDeleted"}).
//		AddRow(farmerData.Id, farmerData.Name, farmerData.District, farmerData.State, farmerData.PhoneNumber, farmerData.IsDeleted)
//	mock.ExpectQuery("^SELECT (.+) FROM farmers where farmerId = \\?$").WithArgs(1).WillReturnRows(rows)
//
//	req, _ := http.NewRequest("GET", "http://localhost/farmers/1", nil)
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected while creating request", err)
//	}
//	w := httptest.NewRecorder()
//
//	router.ServeHTTP(w, req)
//
//	if w.Code != 200 {
//		t.Fatalf("expected status code to be 200, but got: %d", w.Code)
//	}
//
//	responseData := service.Farmer{}
//
//	err = json.Unmarshal(w.Body.Bytes(), &responseData)
//	assert.Equal(t, farmerData, responseData)
//}


//func TestShouldDeleteParticularFarmerDetails(t *testing.T) {
//	router := mux.NewRouter()
//	db, mock, err := sqlmock.New()
//
//	router.HandleFunc("/farmers/{id:[0-9]+}", DeleteFarmer(db)).Methods("PATCH")
//
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//
//	defer db.Close()
//
//	mock.ExpectBegin()
//	mock.ExpectExec("^UPDATE farmers SET isDeleted = 1 WHERE farmerId = \\?$").WithArgs(1).WillReturnResult(sqlmock.NewResult(0,1))
//	mock.ExpectCommit()
//
//
//	req, _ := http.NewRequest("PATCH", "http://localhost/farmers/1", nil)
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected while creating request", err)
//	}
//	w := httptest.NewRecorder()
//
//	router.ServeHTTP(w, req)
//
//	if w.Code != 200 {
//		t.Fatalf("expected status code to be 200, but got: %d", w.Code)
//	}
//
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expections: %s", err)
//	}
//}