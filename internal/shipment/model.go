package shipment

type Shipment struct {
	ID          string `gorm:"column:ID"          json:"ID"`
	UserID      string `gorm:"column:UserID"      json:"UserID"`
	Description string `gorm:"column:Description" json:"Description"` //including purchase date and products, quantity
	Tracking    string `gorm:"column:Tracking"    json:"Tracking"`
	Comment     string `gorm:"column:Comment"     json:"Comment"`
	Date        string `gorm:"column:Date"        json:"Date"` //date of the creation of tracking number
}
