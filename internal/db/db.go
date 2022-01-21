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

func OpenDB(config configs.Config) (db *gorm.DB) {
	if dialect(config.DBDriver) == "mysql" {
		myDB, _ := gorm.Open(mysql.Open(config.DBSource), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		return myDB
	} else {
		return nil
		//so other database supported so far
	}
}
