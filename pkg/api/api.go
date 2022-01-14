package api

import (
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Shipment struct {
	ID          string `gorm:"column:ID"          json:"ID"`
	UserID      string `gorm:"column:UserID"      json:"UserID"`
	Description string `gorm:"column:Description" json:"Description"` //including purchase date and products, quantity
	Tracking    string `gorm:"column:Tracking"    json:"Tracking"`
	Comment     string `gorm:"column:Comment"     json:"Comment"`
	Date        string `gorm:"column:Date"        json:"Date"` //date of the creation of tracking number
}

type User struct {
	ID         string `gorm:"column:ID"         json:"ID"`
	WeChatID   string `gorm:"column:WeChatID"   json:"WeChatID"`
	WeChatName string `gorm:"column:WeChatName" json:"WeChatName"`
	RandomCode string `gorm:"column:RandomCode" json:"RandomCode"`
}

type Admin struct {
	ID       string `gorm:"column:ID"       json:"ID"`
	Name     string `gorm:"column:Name"     json:"Name"`
	Password string `gorm:"column:Password" json:"Password"`
}

type HttpResponse struct {
	Authorization bool        `json:"Authorization"`
	Error         error       `json:"Error"`
	Message       string      `json:"Message"`
	Result        interface{} `json:"Result"`
	StatusCode    int         `json:"StatusCode"`
}

type Claims struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	jwt.StandardClaims
}

func RandomString(length int) string {
	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func NewResponse(w http.ResponseWriter, authorization bool, error error, message string, result interface{}, statusCode int) {
	var newMessage string
	if error != nil {
		newMessage = strconv.Itoa(statusCode) + ":" + message
	} else {
		newMessage = message
	}
	newResponse := HttpResponse{
		Authorization: authorization,
		Error:         error,
		Message:       newMessage,
		Result:        result,
		StatusCode:    statusCode,
	}
	jsonNewResp, err := json.Marshal(newResponse)
	if err != nil {
		temp := HttpResponse{
			Authorization: authorization,
			Error:         errors.New("encode response error"),
			Message:       "",
			Result:        nil,
			StatusCode:    500,
		}
		temp1, _ := json.Marshal(temp)
		w.WriteHeader(500)
		w.Write(temp1)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(jsonNewResp)
}
