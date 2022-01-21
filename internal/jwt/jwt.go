package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"net/http"
	"notifications/configs"
	"notifications/internal/api"
	myError "notifications/internal/error"
	"time"
)

type claims struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	jwt.StandardClaims
}

func TokenGenerate(w http.ResponseWriter, credentialID string, credentialName string) bool {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &claims{
		ID:   credentialID,
		Name: credentialName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(configs.JwtKey)
	if err != nil {
		api.NewHttpResponse(w, myError.NewError(err, "Failed to get the signed JWT.", 500), "", nil)
		return false
	}
	w.Header().Add("Authorization", tokenString)
	return true
}

func TokenParse(w http.ResponseWriter, r *http.Request) bool {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		api.NewHttpResponse(w, myError.NewError(errors.New("authorization error"), "Unauthorized: missing a token.", 401), "", nil)
		return false
	}
	claims := &claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return configs.JwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			api.NewHttpResponse(w, myError.NewError(err, "Unauthorized: signature is invalid.", 401), "", nil)
			return false
		}
		api.NewHttpResponse(w, myError.NewError(err, "Unauthorized.", 401), "", nil)
		return false
	}
	if !token.Valid {
		api.NewHttpResponse(w, myError.NewError(errors.New("token invalid error"), "Unauthorized: token is invalid.", 401), "", nil)
		return false
	}
	return true
}
