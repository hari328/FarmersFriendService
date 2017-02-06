package api

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"net/http"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"github.com/FarmersFriendService/model"
	"github.com/gorilla/mux"
	"bytes"
)

func dummyFarmerOne() model.Farmer {
	return model.Farmer{Id:1, Name:"harish", District:"belgam", State:"Karnataka", PhoneNumber:8989829802, IsDeleted: 0}
}

func dummyFarmerTwo() model.Farmer {
	return model.Farmer{Id:2, Name:"palli", District:"kundapur", State:"Karnataka", PhoneNumber:9099009900, IsDeleted:
	0}
}

func getFarmers() Farmers {
	mockFarmerData := Farmers{List:make([]model.Farmer, 0)}

	mockFarmerData.List = append(mockFarmerData.List, dummyFarmerOne())
	mockFarmerData.List = append(mockFarmerData.List, dummyFarmerTwo())
	return mockFarmerData
}

func mockDbResponse(farmers Farmers) sqlmock.Rows{
	rows := sqlmock.NewRows([]string{"id", "name", "district", "state", "phoneNumber", "isDeleted"}).
		AddRow(farmers.List[0].Id, farmers.List[0].Name, farmers.List[0].District, farmers.List[0].State, farmers.List[0].PhoneNumber, farmers.List[0].IsDeleted).
		AddRow(farmers.List[1].Id, farmers.List[1].Name, farmers.List[1].District, farmers.List[1].State, farmers.List[1].PhoneNumber, farmers.List[0].IsDeleted)
	return rows
}


func TestShouldGetFarmers(t *testing.T) {
	router := mux.NewRouter()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	router.HandleFunc("/farmers", ListFarmers(db)).Methods("GET")

	mockData := getFarmers()
	rows := mockDbResponse(mockData)
	mock.ExpectQuery("^SELECT (.+) FROM farmers$").WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "http://localhost/farmers", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()

	router.ServeHTTP(w,req)

	if w.Code != 200 {
		t.Fatalf("expected status code to be 200, but got: %d", w.Code)
	}

	responseData := Farmers{List:make([]model.Farmer,0)}
	err = json.Unmarshal(w.Body.Bytes(), &responseData)
	assert.Equal(t, mockData, responseData)
}

func TestShouldAddFarmer(t *testing.T) {
	router := mux.NewRouter()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	router.HandleFunc("/farmers", AddFarmer(db)).Methods("POST")

	farmerData := dummyFarmerOne()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO farmers").WithArgs(farmerData.Name, farmerData.District, farmerData.State, farmerData.PhoneNumber, farmerData.IsDeleted).WillReturnResult(sqlmock.NewResult(1,1))
	mock.ExpectCommit()

	farmerJson, err := json.Marshal(farmerData)
	if err != nil {
		t.Fatalf("json marshall failed")
	}
	req, _ := http.NewRequest("POST", "http://localhost/farmers", bytes.NewBuffer([]byte(farmerJson)))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected status code to be 200, but got: %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestShouldGetParticularFarmerDetails(t *testing.T) {
	router := mux.NewRouter()
	db, mock, err := sqlmock.New()

	router.HandleFunc("/farmers/{id:[0-9]+}", GetFarmer(db)).Methods("GET")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()
	farmerData := dummyFarmerOne()

	rows := sqlmock.NewRows([]string{"id", "name", "district", "state", "phoneNumber", "isDeleted"}).
		AddRow(farmerData.Id, farmerData.Name, farmerData.District, farmerData.State, farmerData.PhoneNumber, farmerData.IsDeleted)
	mock.ExpectQuery("^SELECT (.+) FROM farmers where farmerId = \\?$").WithArgs(1).WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "http://localhost/farmers/1", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("expected status code to be 200, but got: %d", w.Code)
	}

	responseData := model.Farmer{}

	err = json.Unmarshal(w.Body.Bytes(), &responseData)
	assert.Equal(t, farmerData, responseData)
}
//
//func TestShouldDeleteParticularFarmerDetails(t *testing.T) {
//	router := mux.NewRouter()
//	db, mock, err := sqlmock.New()
//
//	router.HandleFunc("/farmers/{id:[0-9]+}", DeleteFarmer(db)).Methods("DELETE")
//
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//
//	defer db.Close()
//
//	mock.ExpectBegin()
//	mock.ExpectExec("^DELETE from farmers WHERE farmerId = \\?$").WithArgs(1).WillReturnResult(sqlmock.NewResult(0,1))
//	mock.ExpectCommit()
//
//
//	req, _ := http.NewRequest("DELETE", "http://localhost/farmers/1", nil)
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
//	assert.True(t, mock.ExpectationsWereMet())
//}