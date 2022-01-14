package admin

import (
	"net/http"
	"notifications/pkg/api"
	"notifications/pkg/db"
	"notifications/pkg/jwt"
)

func CreateShipmentAdmin(w http.ResponseWriter, r *http.Request) {
	tkn := jwt.TokenParse(r)
	if tkn.Authorization == true {
		myDB := db.OpenDB(w)
		shipment := api.Shipment{
			ID:          r.FormValue("ID"),
			UserID:      r.FormValue("UserID"),
			Description: r.FormValue("Description"),
			Tracking:    r.FormValue("Tracking"),
			Comment:     r.FormValue("Comment"),
			Date:        r.FormValue("Date"),
		}
		result := myDB.Create(&shipment)
		if result.Error != nil {
			api.NewResponse(w, tkn.Authorization, result.Error, "Database query error.", nil, 404)
			return
		}
		var message = "A new shipment with ID " + shipment.ID + " has been created successfully"
		api.NewResponse(w, tkn.Authorization, nil, message, shipment, 201)
	} else {
		api.NewResponse(w, tkn.Authorization, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}
