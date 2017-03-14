package repository

import (
	"fmt"
	"database/sql"
	"github.com/FarmersFriendService/model"
	"github.com/gocraft/dbr"
)

const (
	IsDeletedColumn          = "isdeleted"
	farmerIdColumn           = "farmerId"
	farmerTable              = "farmers"
	phoneNumberColumn        = "phoneNumber"
	stateColumn              = "state"
	districtColumn           = "district"
	nameColumn               = "name"
	SelectAllFarmersCriteria = -1
)

type FarmerRepository interface {
	ListFarmers() ([]model.Farmer, error)
	AddFarmer(farmerJson []byte) error
	GetFarmer(id int) (model.Farmer, string)
	DeleteFarmer(id int) error
}

type farmerRepository struct {
	Db *dbr.Connection
}

func New(db *dbr.Connection) FarmerRepository {
	return &farmerRepository{Db: db}
}

func (repo farmerRepository) DeleteFarmer(farmerId int) error {
	session := repo.Db.NewSession(nil)
	result, err := session.Update(farmerTable).Set(IsDeletedColumn, 1).
		Where(dbr.Eq(farmerIdColumn, farmerId)).Exec()
	
	return checkResultOnDbModification(err, result, "DeleteFarmer")
}

func (repo farmerRepository) AddFarmer(farmerJson []byte) error {
	farmer, err := model.Unmarshal(farmerJson)
	if err != nil {
		return err
	}
	
	session := repo.Db.NewSession(nil)
	result, err := session.InsertInto(farmerTable).
		Columns(nameColumn, districtColumn, stateColumn, phoneNumberColumn, IsDeletedColumn).
		Record(farmer).Pair(IsDeletedColumn, 0).Exec()
	
	return checkResultOnDbModification(err, result, "AddFarmer")
}

func (repo farmerRepository) ListFarmers() ([]model.Farmer, error) {
	return repo.getFarmers(SelectAllFarmersCriteria)
}

func (repo farmerRepository) GetFarmer(id int) (model.Farmer, string) {
	farmers, err := repo.getFarmers(id)
	if err != nil {
		return model.Farmer{}, err.Error()
	}
	return farmers[0], ""
}

func (repo farmerRepository) getFarmers(farmerId int) ([]model.Farmer, error) {
	farmers := make([]model.Farmer, 0)
	session := repo.Db.NewSession(nil)
	
	selectAllFarmers := session.Select(farmerIdColumn, nameColumn, districtColumn, stateColumn, phoneNumberColumn, IsDeletedColumn).From(farmerTable)
	
	if farmerId != -1 {
		_, err := selectAllFarmers.Where(dbr.Eq(farmerIdColumn, farmerId)).Load(&farmers)
		return farmers, err
	}
	_, err := selectAllFarmers.Load(&farmers)
	return farmers, err
}

func checkResultOnDbModification(err error ,result sql.Result, methodName string) error{
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0{
		return fmt.Errorf("failed in method: %s", methodName)
	}
	return nil
}

