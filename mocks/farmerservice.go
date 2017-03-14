package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/FarmersFriendService/model"
)

type MockFarmerService struct {
	mock.Mock
}

func (m *MockFarmerService) ListFarmers()([]model.Farmer, error) {
	stubArgs := m.Called()
	
	var returnValue []model.Farmer
	returnError := stubArgs.Error(1)
	
	if  returnError == nil {
		returnValue = stubArgs.Get(0).([]model.Farmer)
	}
	return  returnValue, returnError
}


func (m *MockFarmerService) GetFarmer(id int)(model.Farmer, string) {
	stubArgs := m.Called(id)
	returnError := stubArgs.String(1)
	
	var returnValue model.Farmer
	if  returnError == "" {
		returnValue = stubArgs.Get(0).(model.Farmer)
	}
	return  returnValue, returnError
}


func (m *MockFarmerService) AddFarmer(farmer []byte) error {
	args := m.Called(farmer)
	return args.Error(0)
}

func (m *MockFarmerService) DeleteFarmer(id int) error {
	stubArgs := m.Called(id)
	return stubArgs.Error(0)
}