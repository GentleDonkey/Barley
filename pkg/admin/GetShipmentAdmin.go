package admin

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"notifications/pkg/api"
	"notifications/pkg/db"
	"notifications/pkg/jwt"
	"strconv"
)

func GetShipmentAdmin(w http.ResponseWriter, r *http.Request) {
	tkn := jwt.TokenParse(r)
	if tkn.Authorization == true {
		myDB := db.OpenDB()
		shipmentID := mux.Vars(r)["id"]
		_, err := strconv.Atoi(shipmentID)
		if shipmentID == "" || err != nil {
			api.NewResponse(w, tkn.Authorization, err, "Invalid id.", nil, 400)
			return
		}
		cursor, err := myDB.Query("SELECT * FROM shipment WHERE id=?", shipmentID)
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
		if result == nil {
			err := errors.New("none result error")
			api.NewResponse(w, tkn.Authorization, err, "Database does not found any result.", nil, 404)
			return
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
