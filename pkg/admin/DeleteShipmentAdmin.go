package admin

import (
	"github.com/gorilla/mux"
	"net/http"
	"notifications/pkg/api"
	"notifications/pkg/db"
	"notifications/pkg/jwt"
	"strconv"
)

func DeleteShipmentAdmin(w http.ResponseWriter, r *http.Request) {
	tkn := jwt.TokenParse(r)
	if tkn.Authorization == true {
		myDB := db.OpenDB(w)
		shipmentID := mux.Vars(r)["id"]
		_, err := strconv.Atoi(shipmentID)
		if shipmentID == "" || err != nil {
			api.NewResponse(w, tkn.Authorization, err, "Invalid id.", nil, 400)
			return
		}
		result := myDB.Delete(&api.Shipment{}, shipmentID)
		if result.Error != nil {
			api.NewResponse(w, tkn.Authorization, result.Error, "Database query error.", nil, 404)
			return
		}
		var message = "The shipment with ID " + shipmentID + " has been deleted successfully"
		api.NewResponse(w, tkn.Authorization, err, message, nil, 200)
	} else {
		api.NewResponse(w, tkn.Authorization, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}
