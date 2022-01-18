package admin

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type APIAdminRepo interface {
	Login(ta Admin) (*Admin, error, string, int)
}
type adminRepo struct {
	db *gorm.DB
}

func NewAdminRepo(db *gorm.DB) *adminRepo {
	return &adminRepo{
		db,
	}
}

func (a *adminRepo) Login(ta Admin) (*Admin, error, string, int) {
	storedAdmin := &Admin{}
	result := a.db.Raw("SELECT * FROM `admin` WHERE Name=?", ta.Name).Scan(&storedAdmin)
	if result.Error != nil {
		return nil, result.Error, "Database query error.", 404
	}
	err := bcrypt.CompareHashAndPassword([]byte(storedAdmin.Password), []byte(ta.Password))
	if err != nil {
		return nil, err, "Admin name does not match with password.", 401
	}
	newMessage := "Admin " + storedAdmin.Name + " login successfully"
	return storedAdmin, nil, newMessage, 200
}
