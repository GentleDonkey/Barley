package user

import (
	"gorm.io/gorm"
)

type APIUserRepo interface {
	Create(u User) (error, string, int)
	FindAll() ([]User, error, string, int)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepo {
	return &userRepo{
		db,
	}
}

func (u *userRepo) Create(user User) (error, string, int) {
	result := u.db.Create(&user)
	if result.Error != nil {
		return result.Error, "Database query error.", 404
	}
	newMessage := "A new user with ID " + user.ID + " has been created successfully"
	return nil, newMessage, 201
}

func (u *userRepo) FindAll() ([]User, error, string, int) {
	var result []User
	rows, err := u.db.Raw("SELECT * FROM user ORDER BY id DESC").Rows()
	if err != nil {
		return nil, err, "Database query error.", 404
	}
	for rows.Next() {
		err := u.db.ScanRows(rows, &result)
		if err != nil {
			return nil, err, "Database query error.", 404
		}
	}
	newMessage := "All user have been found successfully"
	return result, nil, newMessage, 200
}
