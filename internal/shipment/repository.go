package shipment

import (
	"gorm.io/gorm"
)

type APIShipmentRepo interface {
	Create(s Shipment) (error, string, int)
	FindAll() ([]Shipment, error, string, int)
	FindOne(id int) (Shipment, error, string, int)
	Update(s Shipment) (error, string, int)
	Delete(id int) (error, string, int)
}

type shipmentRepo struct {
	db *gorm.DB
}

func NewShipmentRepo(db *gorm.DB) *shipmentRepo {
	return &shipmentRepo{
		db,
	}
}

func (s *shipmentRepo) Create(shipment Shipment) (error, string, int) {
	result := s.db.Create(&shipment)
	if result.Error != nil {
		return result.Error, "Database query error.", 404
	} else {
		newMessage := "A new shipment with ID " + shipment.ID + " has been created successfully"
		return nil, newMessage, 201
	}
}

func (s *shipmentRepo) FindAll() ([]Shipment, error, string, int) {
	var result []Shipment
	rows, err := s.db.Raw("SELECT * FROM shipment ORDER BY id DESC").Rows()
	if err != nil {
		return nil, err, "Database query error.", 404
	}
	for rows.Next() {
		err := s.db.ScanRows(rows, &result)
		if err != nil {
			return nil, err, "Database query error.", 404
		}
	}
	newMessage := "All shipment have been found successfully"
	return result, nil, newMessage, 200
}

func (s *shipmentRepo) FindOne(shipmentID string) (Shipment, error, string, int) {
	var result Shipment
	s.db.Raw("SELECT * FROM shipment WHERE id=?", shipmentID).Scan(&result)
	newMessage := "The shipment with ID " + result.ID + " has been found successfully"
	return result, nil, newMessage, 200
}

func (s *shipmentRepo) Update(shipment Shipment) (error, string, int) {
	result := s.db.Save(&shipment)
	if result.Error != nil {
		return result.Error, "Database query error.", 404
	}
	newMessage := "The shipment with ID " + shipment.ID + " has been updated successfully"
	return nil, newMessage, 200
}

func (s *shipmentRepo) Delete(shipmentID string) (error, string, int) {
	result := s.db.Delete(&Shipment{}, shipmentID)
	if result.Error != nil {
		return result.Error, "Database query error.", 404
	}
	newMessage := "The shipment with ID " + shipmentID + " has been deleted successfully"
	return nil, newMessage, 200
}
