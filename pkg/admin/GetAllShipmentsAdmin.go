package admin

import (
	"database/sql"
	"net/http"
	"notifications/pkg/api"
	"notifications/pkg/db"
	"notifications/pkg/jwt"
)

func GetAllShipmentsAdmin(w http.ResponseWriter, r *http.Request) {
	tkn := jwt.TokenParse(r)
	if tkn.Authorization == true {
		myDB := db.OpenDB()
		cursor, err := myDB.Query("SELECT * FROM shipment ORDER BY id DESC")
		if err != nil {
			api.NewResponse(w, tkn.Authorization, err, "Invalid SQL query.", nil, 404)
			return
		}
		var result []api.Shipment
		for cursor.Next() {
			var id, userID, description, tracking, comment, date string
			err = cursor.Scan(&id, &userID, &description, &tracking, &comment, &date)
			if err != nil {
				api.NewResponse(w, tkn.Authorization, err, "Database query error.", nil, 404)
				return
			}
			result = append(result, api.Shipment{ID: id, UserID: userID, Description: description, Tracking: tracking, Comment: comment, Date: date})
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
