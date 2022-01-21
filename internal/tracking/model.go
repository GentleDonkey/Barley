package tracking

type Shipment struct {
	ID          string `gorm:"column:ID"          json:"id"`
	UserID      string `gorm:"column:UserID"      json:"userid"`
	Description string `gorm:"column:Description" json:"description"` //including purchase date and products, quantity
	Tracking    string `gorm:"column:Tracking"    json:"tracking"`
	Comment     string `gorm:"column:Comment"     json:"comment"`
	Date        string `gorm:"column:Date"        json:"date"` //date of the creation of tracking number
}
