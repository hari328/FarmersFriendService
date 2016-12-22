package api

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"net/http"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)


func getFarmers() Farmers {
	mockFarmerData := Farmers{List:make([]Farmer,0)}

	mockFarmerData.List = append(mockFarmerData.List, Farmer{Id:1, Name:"harish", District:"belgam", State:"Karnataka", PhoneNumber:8989829802})
	mockFarmerData.List = append(mockFarmerData.List, Farmer{Id:2, Name:"palli", District:"kundapur", State:"Karnataka", PhoneNumber:9099009900})
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

	farmers := Farmers{List:make([]Farmer,0)}
	err = json.Unmarshal(w.Body.Bytes(), &farmers)
	assert.Equal(t, mockData,farmers)
}

func TestShouldAddFarmer(t *testing.T) {

}