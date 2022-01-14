package admin

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"notifications/pkg/api"
	"notifications/pkg/db"
	"notifications/pkg/jwt"
)

func LoginAdmin(w http.ResponseWriter, r *http.Request) {
	myDB := db.OpenDB()
	thisAdmin := &api.Admin{
		ID:       "",
		Name:     r.FormValue("Name"),
		Password: r.FormValue("Password"),
	}
	cursor, err := myDB.Query("SELECT * FROM `admin` WHERE Name=?", thisAdmin.Name)
	if err != nil {
		api.NewResponse(w, false, err, "Invalid SQL query.", nil, 404)
		return
	}
	storedAdmin := &api.Admin{}
	for cursor.Next() {
		err = cursor.Scan(&storedAdmin.ID, &storedAdmin.Name, &storedAdmin.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				api.NewResponse(w, false, err, "Admin does not exist.", nil, 401)
				return
			}
			api.NewResponse(w, false, err, "Database query error.", nil, 404)
			return
		}
	}
	if err = bcrypt.CompareHashAndPassword([]byte(storedAdmin.Password), []byte(thisAdmin.Password)); err != nil {
		api.NewResponse(w, false, err, "Admin name does not match with password.", nil, 401)
		return
	}
	defer func(myDB *sql.DB) {
		err := myDB.Close()
		if err != nil {
			api.NewResponse(w, false, err, "Database closing error.", nil, 404)
			return
		}
	}(myDB)
	tkn := jwt.TokenGenerate(w, storedAdmin.ID, storedAdmin.Name)
	if tkn.Authorization == true {
		api.NewResponse(w, true, err, "", storedAdmin, 200)
	} else {
		api.NewResponse(w, tkn.Authorization, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}
