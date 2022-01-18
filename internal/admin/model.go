package admin

type Admin struct {
	ID       string `gorm:"column:ID"       json:"ID"`
	Name     string `gorm:"column:Name"     json:"Name"`
	Password string `gorm:"column:Password" json:"Password"`
}
