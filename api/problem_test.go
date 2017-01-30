package api

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"

	"github.com/FarmersFriendService/model"
	"time"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"encoding/json"
)

func dummyProblemOne() model.Problem {
	return model.Problem{Id:1, FarmerId:1, ProblemDesc:	"not enough water supply for paddy crop", PostedDate: time.Now(), IsSolved: false }
}

func dummyProblemTwo() model.Problem {
	return model.Problem{Id:2, FarmerId:3, ProblemDesc:	"not enough water supply for paddy crop", PostedDate: time.Now(), IsSolved: true }
}

func getProblems() []model.Problem {
	mockFarmerData := make([]model.Problem,0)

	mockFarmerData = append(mockFarmerData, dummyProblemOne())
	mockFarmerData = append(mockFarmerData, dummyProblemTwo())
	return mockFarmerData
}


func TestShouldGetAllProblems(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockProblems := getProblems()

	rows := sqlmock.NewRows([]string{"id", "name", "district", "state", "phoneNumber"}).
		AddRow(mockProblems[0].Id, mockProblems[0].FarmerId, mockProblems[0].ProblemDesc, mockProblems[0].PostedDate, mockProblems[0].IsSolved).
		AddRow(mockProblems[1].Id, mockProblems[1].FarmerId, mockProblems[1].ProblemDesc, mockProblems[1].PostedDate, mockProblems[1].IsSolved)

	mock.ExpectQuery("^SELECT (.+) FROM problems$").WillReturnRows(rows)

	req, _ := http.NewRequest("GET", "http://localhost/problems", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()

	app := &Api{db}
	app.ListProblems(w, req)

	if w.Code != 200 {
		t.Fatalf("expected status code to be 200, but got: %d", w.Code)
	}

	responseData := make([]model.Problem,0)
	err = json.Unmarshal(w.Body.Bytes(), &responseData)
	assert.Equal(t, mockProblems, responseData)
}