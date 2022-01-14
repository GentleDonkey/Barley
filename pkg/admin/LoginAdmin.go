package admin

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"notifications/pkg/api"
	"notifications/pkg/db"
	"notifications/pkg/jwt"
)

func LoginAdmin(w http.ResponseWriter, r *http.Request) {
	myDB := db.OpenDB(w)
	thisAdmin := &api.Admin{
		ID:       "",
		Name:     r.FormValue("Name"),
		Password: r.FormValue("Password"),
	}
	storedAdmin := &api.Admin{}
	result := myDB.Raw("SELECT * FROM `admin` WHERE Name=?", thisAdmin.Name).Scan(&storedAdmin)
	if result.Error != nil {
		api.NewResponse(w, false, result.Error, "Database query error.", nil, 404)
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(storedAdmin.Password), []byte(thisAdmin.Password))
	if err != nil {
		api.NewResponse(w, false, err, "Admin name does not match with password.", nil, 401)
		return
	}
	tkn := jwt.TokenGenerate(w, storedAdmin.ID, storedAdmin.Name)
	if tkn.Authorization == true {
		api.NewResponse(w, true, err, "", thisAdmin, 200)
	} else {
		api.NewResponse(w, tkn.Authorization, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}
