package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/FarmersFriendService/model"
)

type MockFarmerService struct {
	mock.Mock
}

func (m *MockFarmerService) ListFarmers()([]model.Farmer, string) {
	stubArgs := m.Called()
	
	var returnValue []model.Farmer
	returnError := stubArgs.String(1)
	
	if  returnError == "" {
		returnValue = stubArgs.Get(0).([]model.Farmer)
	}
	
	return  returnValue, returnError
}

func (m *MockFarmerService) AddFarmer(farmer []byte)(bool, string) {
	args := m.Called(farmer)
	return args.Bool(0), args.String(1)
}
