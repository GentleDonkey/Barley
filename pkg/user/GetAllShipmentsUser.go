package user

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"notifications/pkg/api"
	"notifications/pkg/db"
)

func GetAllShipmentsUser(w http.ResponseWriter, r *http.Request) {
	myDB := db.OpenDB(w)
	code := mux.Vars(r)["code"]
	if code == "" {
		api.NewResponse(w, false, errors.New("empty code"), "Empty code parameter.", nil, 400)
		return
	}
	var result []api.Shipment
	rows, err := myDB.Raw("SELECT shipment.* FROM `shipment` JOIN `user` ON (user.ID=shipment.UserID AND user.RandomCode=?)", code).Rows()
	if err != nil {
		api.NewResponse(w, false, err, "Database query error.", nil, 404)
	}
	for rows.Next() {
		err := myDB.ScanRows(rows, &result)
		if err != nil {
			api.NewResponse(w, false, err, "Database query error.", nil, 404)
			return
		}
	}
	api.NewResponse(w, false, err, "", result, 200)
}
