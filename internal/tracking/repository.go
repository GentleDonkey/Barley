package tracking

import (
	"gorm.io/gorm"
	myError "notifications/internal/error"
)

type APITrackingRepo interface {
	FindAll(code string) ([]Shipment, *myError.MyError)
}

type trackingRepo struct {
	db *gorm.DB
}

func NewTrackingRepo(db *gorm.DB) *trackingRepo {
	return &trackingRepo{
		db,
	}
}

func (t *trackingRepo) FindAll(code string) ([]Shipment, *myError.MyError) {
	var result []Shipment
	rows, err := t.db.Raw("SELECT shipment.* FROM `shipment` JOIN `user` ON (user.ID=shipment.UserID AND user.RandomCode=?)", code).Rows()
	if err != nil {
		return nil, myError.NewError(err, "Database query error.", 404)
	}
	for rows.Next() {
		err := t.db.ScanRows(rows, &result)
		if err != nil {
			return nil, myError.NewError(err, "Database query error.", 404)
		}
	}
	return result, nil
}
