package server

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"notifications/configs"
	"notifications/internal/admin"
	"notifications/internal/db"
	"notifications/internal/middleware"
	"notifications/internal/shipment"
	"notifications/internal/tracking"
	"notifications/internal/user"
)

func SetServer(config configs.Config) (r *mux.Router) {

	r = mux.NewRouter().StrictSlash(true)
	myDB := db.OpenDB(config)
	s := r.PathPrefix(config.Version).Subrouter()
	ar := admin.NewAdminRepo(myDB)
	admin.RegisterRoute(ar, s)
	// set admin and add middleware
	adminRouter := s.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.AuthorizationMiddleware)
	sr := shipment.NewShipmentRepo(myDB)
	shipment.RegisterRoute(sr, adminRouter)
	ur := user.NewUserRepo(myDB)
	user.RegisterRoute(ur, adminRouter)
	// set user
	userRouter := s.PathPrefix("/users").Subrouter()
	tr := tracking.NewTrackingRepo(myDB)
	tracking.RegisterRoute(tr, userRouter)
	return r
}
