package admin

import (
	"database/sql"
	"net/http"
	"notifications/pkg/api"
	"notifications/pkg/db"
	"notifications/pkg/jwt"
)

func CreateShipmentAdmin(w http.ResponseWriter, r *http.Request) {
	tkn := jwt.TokenParse(r)
	if tkn.Authorization == true {
		myDB := db.OpenDB()
		cursor, err := myDB.Prepare("INSERT INTO shipment(id, UserID, Description, Tracking, Comment, Date) VALUES(?,?,?,?,?,?)")
		if err != nil {
			api.NewResponse(w, tkn.Authorization, err, "Invalid SQL query.", nil, 404)
			return
		}
		_, err = cursor.Exec(r.FormValue("id"), r.FormValue("UserID"), r.FormValue("Description"), r.FormValue("Tracking"), r.FormValue("Comment"), r.FormValue("Date"))
		if err != nil {
			api.NewResponse(w, tkn.Authorization, err, "Database query error.", nil, 404)
			return
		}
		result := &api.Shipment{
			ID:          r.FormValue("id"),
			UserID:      r.FormValue("UserID"),
			Description: r.FormValue("Description"),
			Tracking:    r.FormValue("Tracking"),
			Comment:     r.FormValue("Comment"),
			Date:        r.FormValue("Date"),
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
