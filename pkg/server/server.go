package server

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"notifications/configs"
	"notifications/internal/admin"
	"notifications/internal/shipment"
	"notifications/internal/tracking"
	"notifications/internal/user"
	"notifications/pkg/db"
)

func SetServer() (r *mux.Router) {
	r = mux.NewRouter().StrictSlash(true)
	myDB := db.OpenDB()
	// set admin
	s := r.PathPrefix(configs.Version).Subrouter()
	adminRouter := s.PathPrefix("/admin").Subrouter()
	ar := admin.NewAdminRepo(myDB)
	admin.RegisterRoute(ar, adminRouter)
	sr := shipment.NewShipmentRepo(myDB)
	shipment.RegisterRoute(sr, adminRouter)
	ur := user.NewUserRepo(myDB)
	user.RegisterRoute(ur, adminRouter)
	// set user
	userRouter := s.PathPrefix("/user").Subrouter()
	tr := tracking.NewTrackingRepo(myDB)
	tracking.RegisterRoute(tr, userRouter)
	return r
}
