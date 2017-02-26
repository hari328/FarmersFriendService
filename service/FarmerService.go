package service

import (
	"fmt"
	"database/sql"
	"github.com/FarmersFriendService/model"
)

type FarmerService interface {
	ListFarmers() ([]model.Farmer, string)
	AddFarmer(farmerJson []byte) (bool, string)
	GetFarmer(id int) (model.Farmer, string)
	DeleteFarmer(id int) error
}

type farmerService struct {
	Db *sql.DB
}

func New(db *sql.DB) FarmerService {
	return &farmerService{Db: db}
}

func (service *farmerService) DeleteFarmer(id int) error {
	fmt.Println(service, id)
	return fmt.Errorf("somethign went wrong")
}

func (service *farmerService) ListFarmers() ([]model.Farmer, string) {
	rows, err := service.Db.Query("SELECT * FROM farmers WHERE isDeleted = 0")
	defer closeSqlRows(err, rows)
	if err != nil {
		return nil, fmt.Sprintf("unable to retrive Farmers data from db: %s", err.Error())
	}
	return getFarmersFromRows(rows)
}

func (service *farmerService) AddFarmer(farmerJson []byte) (bool, string) {
	farmer, err := model.Unmarshal(farmerJson)
	if err != nil {
		return false, err.Error()
	}
	transaction, err := service.Db.Begin()
	if err != nil {
		return false, err.Error()
	}
	defer closeDbTransaction(err, transaction)
	isDeleted := 0
	result, err := transaction.Exec("INSERT INTO farmers(name, district, state, phoneNumber, isDeleted) VALUES (?, ?, ?, ?, ?)", farmer.Name, farmer.District, farmer.State, farmer.PhoneNumber, isDeleted)
	return isDbTransactionSuccessful(result, err)
}

func (service *farmerService) GetFarmer(id int) (model.Farmer, string) {
	rows, err := service.Db.Query("SELECT * FROM farmers where farmerId = ?", id)
	if err != nil {
		return model.Farmer{}, fmt.Sprintf("unable to retrive Farmers data from db: %s", err.Error())
	}
	defer closeSqlRows(err, rows)
	
	res, er := getFarmersFromRows(rows)
	return res[0], er
}

func isDbTransactionSuccessful(result sql.Result, err error) (bool, string) {
	if err != nil {
		return false, err.Error()
	}
	if val, err := result.RowsAffected(); val != 1 || err != nil {
		return false, err.Error()
	}
	return true, ""
}

func getFarmersFromRows(rows *sql.Rows) ([]model.Farmer, string) {
	var farmers []model.Farmer
	var farmer model.Farmer
	
	for rows.Next() {
		err := rows.Scan(&farmer.Id, &farmer.Name, &farmer.District, &farmer.State, &farmer.PhoneNumber, &farmer.IsDeleted)
		if err != nil {
			return nil, fmt.Sprintf("unable to retrive Farmers data from db: %s", err)
		}
		farmers = append(farmers, farmer)
	}
	return farmers, ""
}

func closeSqlRows(err error, rows *sql.Rows) {
	if err == nil {
		rows.Close()
	}
}

func closeDbTransaction(err error, transaction *sql.Tx) {
	switch err {
	case nil:
		err = transaction.Commit()
	default:
		transaction.Rollback()
	}
}
