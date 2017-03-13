package service

import (
	"fmt"
	"database/sql"
	"github.com/FarmersFriendService/model"
	"github.com/gocraft/dbr"
)

type FarmerService interface {
	ListFarmers() ([]model.Farmer, string)
	AddFarmer(farmerJson []byte) error
	GetFarmer(id int) (model.Farmer, string)
	DeleteFarmer(id int) error
}

type farmerService struct {
	Db *dbr.Connection
}

func New(db *dbr.Connection) FarmerService {
	return &farmerService{Db: db}
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

func (service farmerService) DeleteFarmer(farmerId int) error {
	fmt.Println("delete routed")
	session := service.Db.NewSession(nil)
	result, err := session.Update("farmers").Set("isDeleted", 1).
													Where(dbr.Eq("farmerId", farmerId)).Exec()
	
	return checkResultOnDbModification(err, result, "DeleteFarmer")
}

func (service farmerService) getFarmers(farmerId int) ([]model.Farmer, error) {
	session := service.Db.NewSession(nil)
	farmers := make([]model.Farmer, 0)
	
	if farmerId != 0 {
		_, err := session.Select("farmerId", "name", "district", "state", "phoneNumber", "isDeleted").
			From("farmers").Where(dbr.Eq("farmerId", farmerId)).
			Load(&farmers)
		fmt.Println(farmers)
		return farmers, err
	}
	_, err := session.Select("farmerId", "name", "district", "state", "phoneNumber", "isDeleted").
		From("farmers").
		Load(&farmers)
	fmt.Println(farmers)
	return farmers, err
}

func (service farmerService) ListFarmers() ([]model.Farmer, string) {
	farmers, err := service.getFarmers(0)
	if err != nil {
		return nil, err.Error()
	}
	
	return farmers, ""
}

func (service farmerService) GetFarmer(id int) (model.Farmer, string) {
	farmers, err := service.getFarmers(id)
	if err != nil {
		return model.Farmer{}, err.Error()
	}
	return farmers[0], ""
}

func (service farmerService) AddFarmer(farmerJson []byte) error {
	farmer, err := model.Unmarshal(farmerJson)
	if err != nil {
		return err
	}
	
	session := service.Db.NewSession(nil)
	result, err := session.InsertInto("farmers").
		Columns("name", "district", "state", "phoneNumber", "isDeleted").Record(farmer).Pair("isDeleted", 0).Exec()
	 
	return checkResultOnDbModification(err, result, "AddFarmer")
}

