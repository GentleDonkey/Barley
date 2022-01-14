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
	adminRouter.HandleFunc("/shipments/{id}/", admin.GetShipmentAdmin).Methods("GET")
	adminRouter.HandleFunc("/shipments/{id}", admin.DeleteShipmentAdmin).Methods("DELETE")
	adminRouter.HandleFunc("/shipments/{id}", admin.UpdateShipmentAdmin).Methods("PATCH")
	adminRouter.HandleFunc("/users", admin.GetAllUsersAdmin).Methods("GET")
	adminRouter.HandleFunc("/user", admin.CreateUserAdmin).Methods("POST")
	userRouter := s.PathPrefix("/admin").Subrouter()
	userRouter.HandleFunc("/api/v1/user/tracking/{code}", user.GetAllShipmentsUser).Methods("GET")
	return r
}
