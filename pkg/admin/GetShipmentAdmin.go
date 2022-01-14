package admin

import (
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
		myDB := db.OpenDB(w)
		shipmentID := mux.Vars(r)["id"]
		_, err := strconv.Atoi(shipmentID)
		if shipmentID == "" || err != nil {
			api.NewResponse(w, tkn.Authorization, err, "Invalid id.", nil, 400)
			return
		}
		var result []api.Shipment
		rows, err := myDB.Raw("SELECT * FROM shipment WHERE id=?", shipmentID).Rows()
		if err != nil {
			api.NewResponse(w, tkn.Authorization, err, "Database query error.", nil, 404)
		}
		for rows.Next() {
			err := myDB.ScanRows(rows, &result)
			if err != nil {
				api.NewResponse(w, tkn.Authorization, err, "Database query error.", nil, 404)
				return
			}
		}
		api.NewResponse(w, tkn.Authorization, err, "", result, 200)
	} else {
		api.NewResponse(w, tkn.Authorization, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}
