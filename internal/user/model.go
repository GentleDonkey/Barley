package user

type User struct {
	ID         string `gorm:"column:ID"         json:"id"`
	WeChatID   string `gorm:"column:WeChatID"   json:"wechatid"`
	WeChatName string `gorm:"column:WeChatName" json:"wechatname"`
	RandomCode string `gorm:"column:RandomCode" json:"randomcode"`
}
