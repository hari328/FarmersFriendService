package model

import "encoding/json"

type Farmer struct {
	Id          int                `json:"farmerId"`
	Name        string             `json:"name"`
	District    string             `json:"district"`
	State       string             `json:"state"`
	PhoneNumber int64              `json:"phoneNumber"`
	IsDeleted   int                `json:"isDeleted"`
}

func Unmarshal(farmerJson []byte) (Farmer, error) {
	obj := Farmer{}
	if err := json.Unmarshal(farmerJson, &obj); err != nil {
		return Farmer{}, err
	}
	return obj, nil
}
