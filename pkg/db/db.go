package db

import (
	"database/sql"
	"fmt"
	"notifications/configs"
)

func OpenDB() (db *sql.DB) {
	db, err := sql.Open(configs.DbDriver, configs.DbUser+":"+configs.DbPass+"@/"+configs.DbName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	return db
}
