package admin

import (
	"gorm.io/gorm"
	myError "notifications/internal/error"
)

type APIAdminRepo interface {
	Login(ta Admin) (*Admin, *myError.MyError)
}
type adminRepo struct {
	db *gorm.DB
}

func NewAdminRepo(db *gorm.DB) *adminRepo {
	return &adminRepo{
		db,
	}
}

func (a *adminRepo) Login(ta Admin) (*Admin, *myError.MyError) {
	storedAdmin := &Admin{}
	result := a.db.Raw("SELECT * FROM `admin` WHERE Name=?", ta.Name).Scan(&storedAdmin)
	if result.Error != nil {
		return nil, myError.NewError(result.Error, "Database query error.", 404)
	}
	return storedAdmin, nil
}
