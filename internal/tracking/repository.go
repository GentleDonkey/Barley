package tracking

import (
	"gorm.io/gorm"
)

type APITrackingRepo interface {
	FindAll(code string) ([]Shipment, error, string, int)
}

type trackingRepo struct {
	db *gorm.DB
}

func NewTrackingRepo(db *gorm.DB) *trackingRepo {
	return &trackingRepo{
		db,
	}
}

func (t *trackingRepo) FindAll(code string) ([]Shipment, error, string, int) {
	var result []Shipment
	rows, err := t.db.Raw("SELECT shipment.* FROM `shipment` JOIN `user` ON (user.ID=shipment.UserID AND user.RandomCode=?)", code).Rows()
	if err != nil {
		return nil, err, "Database query error.", 404
	}
	for rows.Next() {
		err := t.db.ScanRows(rows, &result)
		if err != nil {
			return nil, err, "Database query error.", 404
		}
	}
	newMessage := "All shipment have been found successfully"
	return result, nil, newMessage, 200
}
