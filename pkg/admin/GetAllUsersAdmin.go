package admin

import (
	"database/sql"
	"net/http"
	"notifications/pkg/api"
	"notifications/pkg/db"
	"notifications/pkg/jwt"
)

func GetAllUsersAdmin(w http.ResponseWriter, r *http.Request) {
	tkn := jwt.TokenParse(r)
	if tkn.Authorization == true {
		myDB := db.OpenDB()
		cursor, err := myDB.Query("SELECT * FROM user ORDER BY id")
		if err != nil {
			api.NewResponse(w, tkn.Authorization, err, "Invalid SQL query.", nil, 404)
			return
		}
		var result []api.User
		for cursor.Next() {
			var id, weChatID, weChatName string
			err = cursor.Scan(&id, &weChatID, &weChatName)
			if err != nil {
				api.NewResponse(w, tkn.Authorization, err, "Database query error.", nil, 404)
				return
			}
			result = append(result, api.User{ID: id, WeChatID: weChatID, WeChatName: weChatName})
		}
		defer func(myDB *sql.DB) {
			err := myDB.Close()
			if err != nil {
				api.NewResponse(w, tkn.Authorization, err, "Database closing error.", nil, 404)
				return
			}
		}(myDB)
		api.NewResponse(w, tkn.Authorization, err, "", result, 200)
	} else {
		api.NewResponse(w, tkn.Authorization, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}
