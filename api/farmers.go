package api

import (
	"database/sql"
)

func AddFarmer(db *sql.DB, farmerName string, farmerPhoneNumber string) (err error) {
	transaction, err := db.Begin()

	_, err = transaction.Exec("INSERT INTO farmers (name, phoneNumber) VALUES (?, ?)", farmerName, farmerPhoneNumber)

		return err
}
