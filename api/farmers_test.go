package api

import (
	"testing"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestShouldAddFarmer(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	farmerName := "girish"
	farmerPhoneNumber := "9012343321"

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO farmers").WithArgs(farmerName, farmerPhoneNumber).WillReturnResult(sqlmock.NewResult(1,1))
	mock.ExpectClose()

	err = addFarmer(db, farmerName, farmerPhoneNumber)
	if err != nil {
		t.Error("unable to add farmer to farmers table")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}


func TestShouldNotAddFarmerIfInputIsNil(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	farmerPhoneNumber := "9012343321"

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO farmers").WithArgs(nil, farmerPhoneNumber).WillReturnError(error("failed to add farmer"))
	mock.ExpectClose()

	err = addFarmer(db, nil, farmerPhoneNumber)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}