package api

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"net/http"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"encoding/json"
	"bytes"
	"github.com/FarmersFriendService/model"
)


func dummyFarmerOne() model.Farmer {
	return model.Farmer{Id:1, Name:"harish", District:"belgam", State:"Karnataka", PhoneNumber:8989829802}
}

func dummyFarmerTwo() model.Farmer {
	return model.Farmer{Id:2, Name:"palli", District:"kundapur", State:"Karnataka", PhoneNumber:9099009900}
}

func getFarmers() Farmers {
	mockFarmerData := Farmers{List:make([]model.Farmer, 0)}

	mockFarmerData.List = append(mockFarmerData.List, dummyFarmerOne())
	mockFarmerData.List = append(mockFarmerData.List, dummyFarmerTwo())
	return mockFarmerData
}

func mockDbResponse(farmers Farmers) sqlmock.Rows{
	rows := sqlmock.NewRows([]string{"id", "name", "district", "state", "phoneNumber"}).
		AddRow(farmers.List[0].Id, farmers.List[0].Name, farmers.List[0].District, farmers.List[0].State, farmers.List[0].PhoneNumber).
		AddRow(farmers.List[1].Id, farmers.List[1].Name, farmers.List[1].District, farmers.List[1].State, farmers.List[1].PhoneNumber)
	return rows
}


func TestShouldGetFarmers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockData := getFarmers()
	rows := mockDbResponse(mockData)
	mock.ExpectQuery("^SELECT (.+) FROM farmers$").WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "http://localhost/farmers", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()

	// now we execute our request
	app := &Api{db}
	app.ListFarmers(w, req)

	if w.Code != 200 {
		t.Fatalf("expected status code to be 200, but got: %d", w.Code)
	}

	responseData := Farmers{List:make([]model.Farmer,0)}
	err = json.Unmarshal(w.Body.Bytes(), &responseData)
	assert.Equal(t, mockData, responseData)
}

func TestShouldAddFarmer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	farmerData := dummyFarmerOne()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO farmers").WithArgs(farmerData.Name, farmerData.District, farmerData.State, farmerData.PhoneNumber).WillReturnResult(sqlmock.NewResult(1,1))
	mock.ExpectCommit()

	farmerJson, err := json.Marshal(farmerData)
	if err != nil {
		t.Fatalf("json marshall failed")
	}
	req, _ := http.NewRequest("POST", "http://localhost/farmers", bytes.NewBuffer([]byte(farmerJson)))
	w := httptest.NewRecorder()

	app := &Api{db}
	app.AddFarmer(w, req)

	if w.Code != 200 {
		t.Fatalf("expected status code to be 200, but got: %d", w.Code)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestShouldGetParticularFarmerDetails(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	farmerData := dummyFarmerOne()

	rows := sqlmock.NewRows([]string{"id", "name", "district", "state", "phoneNumber"}).
		AddRow(farmerData.Id, farmerData.Name, farmerData.District, farmerData.State, farmerData.PhoneNumber)
	mock.ExpectQuery("^SELECT (.+) FROM farmers where farmerId = \\?$").WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "http://localhost/farmers/1", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()

	// now we execute our request
	app := &Api{db}
	app.GetFarmer(w, req)

	if w.Code != 200 {
		t.Fatalf("expected status code to be 200, but got: %d", w.Code)
	}

	responseData := model.Farmer{}

	err = json.Unmarshal(w.Body.Bytes(), &responseData)
	assert.Equal(t, farmerData, responseData)
}