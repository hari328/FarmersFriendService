package api

import "database/sql"

type Api struct {
	Db *sql.DB
}

func NewApi(filename string) Api {
	db, err := sql.Open("sqlite3", filename)
	if err != nil { panic(err) }
	if db == nil { panic("db nil") }

	return Api{
		Db : db,
	}
}
