package api

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
	"net/http"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

func TestShouldGetFarmers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	app := &Api{db}

	req, _ := http.NewRequest("GET", "http://localhost/farmers", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}

	mockFarmerData := Farmers{List:make([]Farmer,0)}
	mockFarmerData.List = append(mockFarmerData.List, Farmer{Id:1, Name:"harish", District:"belgam", State:"Karnataka", PhoneNumber:8989829802})
	mockFarmerData.List = append(mockFarmerData.List, Farmer{Id:2, Name:"palli", District:"kundapur", State:"Karnataka", PhoneNumber:9099009900})

	w := httptest.NewRecorder()
	rows := sqlmock.NewRows([]string{"id", "name", "district", "state", "phoneNumber"}).
		AddRow(mockFarmerData.List[0].Id, mockFarmerData.List[0].Name, mockFarmerData.List[0].District, mockFarmerData.List[0].State, mockFarmerData.List[0].PhoneNumber).
		AddRow(mockFarmerData.List[1].Id, mockFarmerData.List[1].Name, mockFarmerData.List[1].District, mockFarmerData.List[1].State, mockFarmerData.List[1].PhoneNumber)

	mock.ExpectQuery("^SELECT (.+) FROM farmers$").WillReturnRows(rows)

	// now we execute our request
	app.ListFarmers(w, req)

	if w.Code != 200 {
		t.Fatalf("expected status code to be 200, but got: %d", w.Code)
	}

	farmers := Farmers{List:make([]Farmer,0)}
	err = json.Unmarshal(w.Body.Bytes(), &farmers)

	assert.Equal(t,mockFarmerData,farmers)
}
