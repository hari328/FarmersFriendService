package util

import (
	"github.com/FarmersFriendService/model"
)

func DummyFarmerOne() model.Farmer {
	return model.Farmer{Id:1, Name:"harish", District:"belgam", State:"Karnataka", PhoneNumber:8989829802, IsDeleted: 0}
}

func dummyFarmerTwo() model.Farmer {
	return model.Farmer{Id:2, Name:"palli", District:"kundapur", State:"Karnataka", PhoneNumber:9099009900, IsDeleted:
	0}
}

func GetDummyFarmers() []model.Farmer {
	mockFarmerData := make([]model.Farmer, 0)
	mockFarmerData = append(mockFarmerData, DummyFarmerOne())
	mockFarmerData = append(mockFarmerData, dummyFarmerTwo())
	return mockFarmerData
}

type DbError struct {
	Err string
}

func (dbError DbError) Error() string {
	return dbError.Err
}
