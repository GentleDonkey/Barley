package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"net/http"
	"notifications/configs"
	"notifications/pkg/api"
	"time"
)

func TokenGenerate(w http.ResponseWriter, credentialID string, credentialName string) api.HttpResponse {
	newResponse := api.HttpResponse{}
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &api.Claims{
		ID:   credentialID,
		Name: credentialName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(configs.JwtKey)
	if err != nil {
		newResponse.Authorization = false
		newResponse.Error = err
		newResponse.Message = "Failed to get the signed JWT."
		newResponse.Result = nil
		newResponse.StatusCode = 500
		return newResponse
	}
	w.Header().Add("Authorization", tokenString)
	newResponse.Authorization = true
	newResponse.Error = nil
	newResponse.Message = ""
	newResponse.Result = tokenString
	newResponse.StatusCode = 200
	return newResponse
}

func TokenParse(r *http.Request) api.HttpResponse {
	newResponse := api.HttpResponse{}
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		newResponse.Authorization = false
		newResponse.Error = errors.New("authorization not set")
		newResponse.Message = "Unauthorized: missing a token."
		newResponse.Result = nil
		newResponse.StatusCode = 401
		return newResponse
	}
	claims := &api.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return configs.JwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			newResponse.Authorization = false
			newResponse.Error = err
			newResponse.Message = "Unauthorized: signature is invalid."
			newResponse.Result = nil
			newResponse.StatusCode = 401
			return newResponse
		}
		newResponse.Authorization = false
		newResponse.Error = err
		newResponse.Message = "Unauthorized: internal error."
		newResponse.Result = nil
		newResponse.StatusCode = 400
		return newResponse
	}
	if !token.Valid {
		newResponse.Authorization = false
		newResponse.Error = err
		newResponse.Message = "Unauthorized: token is invalid."
		newResponse.Result = nil
		newResponse.StatusCode = 400
		return newResponse
	}
	newResponse.Authorization = true
	newResponse.Error = nil
	newResponse.Message = ""
	newResponse.Result = token
	newResponse.StatusCode = 200
	return newResponse
}
