package admin

import (
	"database/sql"
	"net/http"
	"notifications/pkg/api"
	"notifications/pkg/db"
	"notifications/pkg/jwt"
)

func CreateUserAdmin(w http.ResponseWriter, r *http.Request) {
	tkn := jwt.TokenParse(r)
	if tkn.Authorization == true {
		myDB := db.OpenDB()
		code := api.RandomString(16)
		cursor, err := myDB.Prepare("INSERT INTO user(id, WeChatID, WeChatName, RandomCode) VALUES(?,?,?,?)")
		if err != nil {
			api.NewResponse(w, tkn.Authorization, err, "Invalid SQL query.", nil, 404)
			return
		}
		_, err = cursor.Exec(r.FormValue("id"), r.FormValue("WeChatID"), r.FormValue("WeChatName"), code)
		if err != nil {
			api.NewResponse(w, tkn.Authorization, err, "Database query error.", nil, 404)
			return
		}
		result := &api.User{
			ID:         r.FormValue("id"),
			WeChatID:   r.FormValue("WeChatID"),
			WeChatName: r.FormValue("WeChatName"),
			RandomCode: code,
		}
		defer func(myDB *sql.DB) {
			err := myDB.Close()
			if err != nil {
				api.NewResponse(w, tkn.Authorization, err, "Database closing error.", nil, 404)
				return
			}
		}(myDB)
		api.NewResponse(w, tkn.Authorization, err, "", result, 201)
	} else {
		api.NewResponse(w, tkn.Authorization, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}
