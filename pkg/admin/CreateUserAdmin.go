package admin

import (
	"net/http"
	"notifications/pkg/api"
	"notifications/pkg/db"
	"notifications/pkg/jwt"
)

func CreateUserAdmin(w http.ResponseWriter, r *http.Request) {
	tkn := jwt.TokenParse(r)
	if tkn.Authorization == true {
		myDB := db.OpenDB(w)
		code := api.RandomString(16)
		user := &api.User{
			ID:         r.FormValue("ID"),
			WeChatID:   r.FormValue("WeChatID"),
			WeChatName: r.FormValue("WeChatName"),
			RandomCode: code,
		}
		result := myDB.Create(&user)
		if result.Error != nil {
			api.NewResponse(w, tkn.Authorization, result.Error, "Database query error.", nil, 404)
			return
		}
		var message = "A new user with ID " + user.ID + " has been created successfully"
		api.NewResponse(w, tkn.Authorization, nil, message, user, 201)
	} else {
		api.NewResponse(w, tkn.Authorization, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}
