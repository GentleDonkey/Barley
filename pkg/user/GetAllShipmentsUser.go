package user

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"notifications/pkg/api"
	"notifications/pkg/db"
)

func GetAllShipmentsUser(w http.ResponseWriter, r *http.Request) {
	myDB := db.OpenDB()
	code := mux.Vars(r)["code"]
	if code == "" {
		api.NewResponse(w, false, errors.New("empty code"), "Empty code parameter.", nil, 400)
		return
	}
	cursor, err := myDB.Query("SELECT shipment.* FROM shipment JOIN user ON (user.id=shipment.userID AND user.RandomCode=?)", code)
	if err != nil {
		api.NewResponse(w, false, err, "Invalid SQL query.", nil, 404)
		return
	}
	var result []api.Shipment
	for cursor.Next() {
		var id, userID, description, tracking, comment, date string
		err = cursor.Scan(&id, &userID, &description, &tracking, &comment, &date)
		if err != nil {
			api.NewResponse(w, false, err, "Database query error.", nil, 404)
			return
		}
		result = append(result, api.Shipment{ID: id, UserID: userID, Description: description, Tracking: tracking, Comment: comment, Date: date})
	}
	if result == nil {
		err := errors.New("none result error")
		api.NewResponse(w, false, err, "Database does not found any result.", nil, 404)
		return
	}
	defer func(myDB *sql.DB) {
		err := myDB.Close()
		if err != nil {
			api.NewResponse(w, false, err, "Database closing error.", nil, 404)
			return
		}
	}(myDB)
	api.NewResponse(w, false, err, "", result, 200)
}
