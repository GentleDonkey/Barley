package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"notifications/configs"
)

func dialect(DbDriver string) (d string) {
	switch DbDriver {
	case "MySQL":
		return "mysql"
	case "PostgreSQL":
		return "postgres"
	}
	return
}

func OpenDB() (db *gorm.DB) {
	dsnString := configs.DbUser + ":" + configs.DbPass + "@tcp(" + configs.Host + configs.DbPath + ")/" + configs.DbName + "?charset=utf8&parseTime=True&loc=Local"
	if dialect(configs.DbDriver) == "mysql" {
		myDB, _ := gorm.Open(mysql.Open(dsnString), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		// need an error handler here
		return myDB
	} else {
		return
		//so other database supported so far
	}
}
