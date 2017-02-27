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
	
	farmerService := New(db)
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
	
	farmerService := New(db)
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
	
	mock.ExpectQuery("^SELECT (.+) FROM farmers WHERE isDeleted = 0$").
		WillReturnError(fmt.Errorf("db error"))
	
	farmerService := New(db)
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
	
	farmerJson, _ := json.Marshal(util.DummyFarmerOne())
	farmerService := New(db)
	er := farmerService.AddFarmer(farmerJson)
	
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	assert.Nil(t,er)
}

func TestShouldNotAddFarmerOnInvalidInput(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	farmerJson := []byte(`{"a":123, "b":"something"}`)
	farmerService := New(db)
	err = farmerService.AddFarmer(farmerJson)
	assert.NotEmpty(t, err.Error())
}

func TestShouldGetFarmerFromDb(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	
	farmer := util.DummyFarmerOne()
	setUpMockForGetFarmer(mock, farmer)
	
	farmerService := New(db)
	res, _ := farmerService.GetFarmer(farmer.Id)
	
	assert.Equal(t, farmer, res)
}

func TestShouldGetEmptyFarmerIfValueDoNotExistInDb(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	
	farmer := util.DummyFarmerOne()
	setUpMockForGetFarmer(mock, model.Farmer{})
	
	farmerService := New(db)
	res, _ := farmerService.GetFarmer(farmer.Id)
	
	assert.Equal(t, model.Farmer{}, res)
}

func TestShouldGetErrorIfDbErrorOutForGetFarmer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	
	farmer := util.DummyFarmerOne()
	mock.ExpectQuery("^SELECT (.+) FROM farmers where farmerId = \\?$").WithArgs(1).
		WillReturnError(fmt.Errorf("db error"))
	
	farmerService := New(db)
	_, er := farmerService.GetFarmer(farmer.Id)
	
	assert.Equal(t, "unable to retrive Farmers data from db: db error", er)
}

func TestShouldDeleteFarmerFromDb(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	
	setUpMockForDeleteFarmerSuccess(mock)
	
	farmerService := New(db)
	err = farmerService.DeleteFarmer(1)
	
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	assert.Nil(t, err)
}

func TestShouldNotDeleteFarmerIfNotFoundInDb(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	
	setupMockForDeleteFarmerNotFoundInDb(mock)
	
	farmerService := New(db)
	er := farmerService.DeleteFarmer(1)
	
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	assert.Equal(t, "unable to find record", er.Error())
}

func TestShouldNotDeleteFarmerOnDbError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	
	setupMockForDeleteFarmerDbFailure(mock)
	
	farmerService := New(db)
	er := farmerService.DeleteFarmer(1)
	
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
	assert.Equal(t, "db failure", er.Error())
}

func setupMockForDeleteFarmerDbFailure(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE farmers SET isDeleted = 1 WHERE farmerId = \\?$").
		WithArgs(1).WillReturnError(fmt.Errorf("db failure"))
	mock.ExpectCommit()
}

func setupMockForDeleteFarmerNotFoundInDb(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE farmers SET isDeleted = 1 WHERE farmerId = \\?$").
		WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
}

func setUpMockForDeleteFarmerSuccess(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE farmers SET isDeleted = 1 WHERE farmerId = \\?$").
		WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()
}

func setUpMockForGetFarmer(mock sqlmock.Sqlmock, farmer model.Farmer) {
	rows := sqlmock.NewRows([]string{"id", "name", "district", "state", "phoneNumber", "isDeleted"}).
		AddRow(farmer.Id, farmer.Name, farmer.District, farmer.State, farmer.PhoneNumber, farmer.IsDeleted)
	mock.ExpectQuery("^SELECT (.+) FROM farmers where farmerId = \\?$").WithArgs(1).
		WillReturnRows(rows)
}

func setUpMockForListFarmers(mock sqlmock.Sqlmock, farmers []model.Farmer) {
	mock.ExpectQuery("^SELECT (.+) FROM farmers WHERE isDeleted = 0$").
		WillReturnRows(mockDbResponse(farmers))
}

func setUpMockForAddFarmer(farmer model.Farmer, mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO farmers").
		WithArgs(farmer.Name, farmer.District, farmer.State, farmer.PhoneNumber, farmer.IsDeleted).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
}

func mockDbResponse(farmers []model.Farmer) sqlmock.Rows {
	if len(farmers) == 0 {
		return sqlmock.NewRows([]string{"id", "name", "district", "state", "phoneNumber", "isDeleted"})
	}
	
	rows := sqlmock.NewRows([]string{"id", "name", "district", "state", "phoneNumber", "isDeleted"}).
		AddRow(farmers[0].Id, farmers[0].Name, farmers[0].District, farmers[0].State, farmers[0].PhoneNumber, farmers[0].IsDeleted).
		AddRow(farmers[1].Id, farmers[1].Name, farmers[1].District, farmers[1].State, farmers[1].PhoneNumber, farmers[0].IsDeleted)
	return rows
}

//func TestShouldNotAddFarmerOnInvalidInput(t *testing.T) {
//	db, _, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	defer db.Close()
//
//	farmerJson := []byte(`{"a":123, "b":"something"}`)
//	farmerService := New(db)
//	isAdded, er := farmerService.AddFarmer(farmerJson)
//
//	assert.Equal(t,false, isAdded)
//	assert.NotEmpty(t, er)
//}
//
//func TestShouldGetFarmerFromDb(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	defer db.Close()
//
//	farmer := util.DummyFarmerOne()
//	setUpMockForGetFarmer(mock, farmer)
//
//	farmerService := New(db)
//	res, _ := farmerService.GetFarmer(farmer.Id)
//
//	assert.Equal(t, farmer, res)
//}
