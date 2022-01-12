package admin

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
	"notifications/pkg/api"
	"notifications/pkg/db"
	"notifications/pkg/jwt"
	"strconv"
)

func UpdateShipmentAdmin(w http.ResponseWriter, r *http.Request) {
	tkn := jwt.TokenParse(r)
	if tkn.Authorization == true {
		myDB := db.OpenDB()
		shipmentID := mux.Vars(r)["id"]
		_, err := strconv.Atoi(shipmentID)
		if shipmentID == "" || err != nil {
			api.NewResponse(w, tkn.Authorization, err, "Invalid id.", nil, 400)
			return
		}
		cursor, err := myDB.Prepare("UPDATE shipment SET UserID=?, Description=?, Tracking=?, Comment=?, Date=? WHERE id=?")
		if err != nil {
			api.NewResponse(w, tkn.Authorization, err, "Invalid SQL query.", nil, 404)
			return
		}
		_, err = cursor.Exec(r.FormValue("UserID"), r.FormValue("Description"), r.FormValue("Tracking"), r.FormValue("Comment"), r.FormValue("Date"), shipmentID)
		if err != nil {
			api.NewResponse(w, tkn.Authorization, err, "Database query error.", nil, 404)
			return
		}
		result := &api.Shipment{
			ID:          shipmentID,
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
		var message = "The shipment with ID " + shipmentID + " has been updated successfully"
		api.NewResponse(w, tkn.Authorization, err, message, result, 200)
	} else {
		api.NewResponse(w, tkn.Authorization, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}
