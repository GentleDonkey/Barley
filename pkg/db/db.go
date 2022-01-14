package db

import (
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"net/http"
	"notifications/configs"
	"notifications/pkg/api"
)

func dialect(w http.ResponseWriter, DbDriver string) (d string) {
	switch DbDriver {
	case "MySQL":
		return "mysql"
	case "PostgreSQL":
		return "postgres"
	}
	api.NewResponse(w, false, errors.New("database error"), "This database is unsupported.", nil, 500)
	return
}

func OpenDB(w http.ResponseWriter) (db *gorm.DB) {
	dsnString := configs.DbUser + ":" + configs.DbPass + "@tcp(" + configs.Host + configs.DbPath + ")/" + configs.DbName + "?charset=utf8&parseTime=True&loc=Local"
	if dialect(w, configs.DbDriver) == "mysql" {
		myDB, err := gorm.Open(mysql.Open(dsnString), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err != nil {
			api.NewResponse(w, false, err, "Fail to connect database.", nil, 500)
			return
		}
		return myDB
	} else {
		return
		//so other database supported so far
	}
}
