package model

type Farmer struct {
	Id int								`json:"farmerId"`
	Name string						`json:"name"`
	District string				`json:"district"`
	State string					`json:"state"`
	PhoneNumber int64			`json:"phoneNumber"`
}
