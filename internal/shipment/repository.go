package shipment

import (
	"gorm.io/gorm"
	myError "notifications/internal/error"
)

type APIShipmentRepo interface {
	Create(s Shipment) *myError.MyError
	FindAll() ([]Shipment, *myError.MyError)
	FindOne(id int) Shipment
	Update(s Shipment) *myError.MyError
	Delete(id int) *myError.MyError
}

type shipmentRepo struct {
	db *gorm.DB
}

func NewShipmentRepo(db *gorm.DB) *shipmentRepo {
	return &shipmentRepo{
		db,
	}
}

func (s *shipmentRepo) Create(shipment Shipment) *myError.MyError {
	result := s.db.Create(&shipment)
	if result.Error != nil {
		return myError.NewError(result.Error, "Database query error.", 404)
	}
	return nil
}

func (s *shipmentRepo) FindAll() ([]Shipment, *myError.MyError) {
	var result []Shipment
	rows, err := s.db.Raw("SELECT * FROM shipment ORDER BY id DESC").Rows()
	if err != nil {
		return nil, myError.NewError(err, "Database query error.", 404)
	}
	for rows.Next() {
		err := s.db.ScanRows(rows, &result)
		if err != nil {
			return nil, myError.NewError(err, "Database query error.", 404)
		}
	}
	return result, nil
}

func (s *shipmentRepo) FindOne(shipmentID string) Shipment {
	var result Shipment
	s.db.Raw("SELECT * FROM shipment WHERE id=?", shipmentID).Scan(&result)
	return result
}

func (s *shipmentRepo) Update(shipment Shipment) *myError.MyError {
	result := s.db.Save(&shipment)
	if result.Error != nil {
		return myError.NewError(result.Error, "Database query error.", 404)
	}
	return nil
}

func (s *shipmentRepo) Delete(shipmentID string) *myError.MyError {
	result := s.db.Delete(&Shipment{}, shipmentID)
	if result.Error != nil {
		return myError.NewError(result.Error, "Database query error.", 404)
	}
	return nil
}
