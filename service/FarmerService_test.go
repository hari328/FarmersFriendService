package service

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/FarmersFriendService/util"
	"github.com/FarmersFriendService/model"
	"fmt"
	"encoding/json"
)

func TestShouldFetchFarmersWhenDbHasFarmers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	setUpMockForListFarmers(mock, util.GetDummyFarmers())
	
	farmerService := NewFarmerService(db)
	resp, _ := farmerService.ListFarmers()
	
	assert.Equal(t, util.GetDummyFarmers(), resp)
}

func TestShouldNotFetchFarmersWhenDbHasNoFarmers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	
	setUpMockForListFarmers(mock, make([]model.Farmer, 0))
	
	farmerService := NewFarmerService(db)
	resp, _ := farmerService.ListFarmers()
	
	var expected []model.Farmer
	assert.Equal(t, expected, resp)
}

func TestShouldNotFetchFarmersWhenDbErrorOut(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	
	mock.ExpectQuery("^SELECT (.+) FROM farmers WHERE isDeleted = 0$").WillReturnError(fmt.Errorf("db error"))
	
	farmerService := NewFarmerService(db)
	_, errString := farmerService.ListFarmers()
	assert.Equal(t, "unable to retrive Farmers data from db: db error", errString)
}

func TestShouldAddFarmerToDb(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	
	setUpMockForAddFarmer(util.DummyFarmerOne(), mock)
	
	farmerJson, _:= json.Marshal(util.DummyFarmerOne())
	farmerService := NewFarmerService(db)
	isAdded, er := farmerService.AddFarmer(farmerJson)
	
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	assert.Equal(t,true, isAdded)
	assert.Equal(t ,"", er)
}

func setUpMockForListFarmers(mock sqlmock.Sqlmock, farmers []model.Farmer) {
	mock.ExpectQuery("^SELECT (.+) FROM farmers WHERE isDeleted = 0$").WillReturnRows(mockDbResponse(farmers))
}

func setUpMockForAddFarmer(farmer model.Farmer, mock sqlmock.Sqlmock ) {
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO farmers").
		WithArgs(farmer.Name, farmer.District, farmer.State, farmer.PhoneNumber, farmer.IsDeleted).
		WillReturnResult(sqlmock.NewResult(1,1))
	mock.ExpectCommit()
}

func mockDbResponse(farmers []model.Farmer) sqlmock.Rows{
	if len(farmers) == 0 {
		return sqlmock.NewRows([]string{"id", "name", "district", "state", "phoneNumber", "isDeleted"})
	}
	
	rows := sqlmock.NewRows([]string{"id", "name", "district", "state", "phoneNumber", "isDeleted"}).
		AddRow(farmers[0].Id, farmers[0].Name, farmers[0].District, farmers[0].State, farmers[0].PhoneNumber, farmers[0].IsDeleted).
		AddRow(farmers[1].Id, farmers[1].Name, farmers[1].District, farmers[1].State, farmers[1].PhoneNumber, farmers[0].IsDeleted)
	return rows
}

////func TestShouldGetParticularFarmerDetails(t *testing.T) {
////	router := mux.NewRouter()
////	db, mock, err := sqlmock.New()
////
////	router.HandleFunc("/farmers/{id:[0-9]+}", GetFarmer(db)).Methods("GET")
////
////	if err != nil {
////		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
////	}
////
////	defer db.Close()
////	farmerData := dummyFarmerOne()
////
////	rows := sqlmock.NewRows([]string{"id", "name", "district", "state", "phoneNumber", "isDeleted"}).
////		AddRow(farmerData.Id, farmerData.Name, farmerData.District, farmerData.State, farmerData.PhoneNumber, farmerData.IsDeleted)
////	mock.ExpectQuery("^SELECT (.+) FROM farmers where farmerId = \\?$").WithArgs(1).WillReturnRows(rows)
////
////	req, _ := http.NewRequest("GET", "http://localhost/farmers/1", nil)
////	if err != nil {
////		t.Fatalf("an error '%s' was not expected while creating request", err)
////	}
////	w := httptest.NewRecorder()
////
////	router.ServeHTTP(w, req)
////
////	if w.Code != 200 {
////		t.Fatalf("expected status code to be 200, but got: %d", w.Code)
////	}
////
////	responseData := service.Farmer{}
////
////	err = json.Unmarshal(w.Body.Bytes(), &responseData)
////	assert.Equal(t, farmerData, responseData)
////}
//
////func TestShouldDeleteParticularFarmerDetails(t *testing.T) {
////	router := mux.NewRouter()
////	db, mock, err := sqlmock.New()
////
////	router.HandleFunc("/farmers/{id:[0-9]+}", DeleteFarmer(db)).Methods("PATCH")
////
////	if err != nil {
////		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
////	}
////
////	defer db.Close()
////
////	mock.ExpectBegin()
////	mock.ExpectExec("^UPDATE farmers SET isDeleted = 1 WHERE farmerId = \\?$").WithArgs(1).WillReturnResult(sqlmock.NewResult(0,1))
////	mock.ExpectCommit()
////
////
////	req, _ := http.NewRequest("PATCH", "http://localhost/farmers/1", nil)
////	if err != nil {
////		t.Fatalf("an error '%s' was not expected while creating request", err)
////	}
////	w := httptest.NewRecorder()
////
////	router.ServeHTTP(w, req)
////
////	if w.Code != 200 {
////		t.Fatalf("expected status code to be 200, but got: %d", w.Code)
////	}
////
////	if err := mock.ExpectationsWereMet(); err != nil {
////		t.Errorf("there were unfulfilled expections: %s", err)
////	}
////}
