package user

type User struct {
	ID         string `gorm:"column:ID"         json:"ID"`
	WeChatID   string `gorm:"column:WeChatID"   json:"WeChatID"`
	WeChatName string `gorm:"column:WeChatName" json:"WeChatName"`
	RandomCode string `gorm:"column:RandomCode" json:"RandomCode"`
}
