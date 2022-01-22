package admin

type Admin struct {
	ID       string `gorm:"column:ID"       json:"id"`
	Name     string `gorm:"column:Name"     json:"name"`
	Password string `gorm:"column:Password" json:"password"`
}
