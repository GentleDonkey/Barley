package admin

import (
	"net/http"
	"notifications/pkg/api"
	"notifications/pkg/db"
	"notifications/pkg/jwt"
)

func GetAllShipmentsAdmin(w http.ResponseWriter, r *http.Request) {
	tkn := jwt.TokenParse(r)
	if tkn.Authorization == true {
		myDB := db.OpenDB(w)
		var result []api.Shipment
		rows, err := myDB.Raw("SELECT * FROM shipment ORDER BY id DESC").Rows()
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
		api.NewResponse(w, tkn.Authorization, nil, "", result, 201)
	} else {
		api.NewResponse(w, tkn.Authorization, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}
