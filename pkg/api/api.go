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
