package server

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"notifications/configs"
	"notifications/pkg/admin"
	"notifications/pkg/user"
)

func SetServer() (r *mux.Router) {
	r = mux.NewRouter().StrictSlash(true)
	s := r.PathPrefix(configs.Version).Subrouter()
	adminRouter := s.PathPrefix("/admin").Subrouter()
	adminRouter.HandleFunc("/login", admin.LoginAdmin).Methods("POST")
	adminRouter.HandleFunc("/shipments", admin.GetAllShipmentsAdmin).Methods("GET")
	adminRouter.HandleFunc("/shipment", admin.CreateShipmentAdmin).Methods("POST")
	adminRouter.HandleFunc("/shipment/{id}/", admin.GetShipmentAdmin).Methods("GET")
	adminRouter.HandleFunc("/shipment/{id}", admin.DeleteShipmentAdmin).Methods("DELETE")
	adminRouter.HandleFunc("/shipment/{id}", admin.UpdateShipmentAdmin).Methods("PATCH")
	adminRouter.HandleFunc("/users", admin.GetAllUsersAdmin).Methods("GET")
	adminRouter.HandleFunc("/user", admin.CreateUserAdmin).Methods("POST")
	userRouter := s.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/tracking/{code}", user.GetAllShipmentsUser).Methods("GET")
	return r
}
