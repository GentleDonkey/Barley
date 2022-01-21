package user

import (
	"gorm.io/gorm"
	myError "notifications/internal/error"
)

type APIUserRepo interface {
	Create(u User) *myError.MyError
	FindAll() ([]User, *myError.MyError)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepo {
	return &userRepo{
		db,
	}
}

func (u *userRepo) Create(user User) *myError.MyError {
	result := u.db.Create(&user)
	if result.Error != nil {
		return myError.NewError(result.Error, "Database query error.", 404)
	}
	return nil
}

func (u *userRepo) FindAll() ([]User, *myError.MyError) {
	var result []User
	rows, err := u.db.Raw("SELECT * FROM user ORDER BY id DESC").Rows()
	if err != nil {
		return nil, myError.NewError(err, "Database query error.", 404)
	}
	for rows.Next() {
		err := u.db.ScanRows(rows, &result)
		if err != nil {
			return nil, myError.NewError(err, "Database query error.", 404)
		}
	}
	return result, nil
}
