package model

import "encoding/json"

type Farmer struct {
	Id          int                `db:"farmerId"`
	Name        string             `db:"name"`
	District    string             `db:"district"`
	State       string             `db:"state"`
	PhoneNumber int64              `db:"phoneNumber"`
	IsDeleted   int                `db:"isdeleted"`
}

func Unmarshal(farmerJson []byte) (Farmer, error) {
	obj := Farmer{}
	if err := json.Unmarshal(farmerJson, &obj); err != nil {
		return Farmer{}, err
	}
	return obj, nil
}
