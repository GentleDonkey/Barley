package admin

import (
	"github.com/gorilla/mux"
	"net/http"
	"notifications/pkg/api"
	"notifications/pkg/jwt"
)

type APIAdminHandler struct {
	repo *adminRepo
}

func NewAdminHandler(ar *adminRepo) *APIAdminHandler {
	return &APIAdminHandler{
		ar,
	}
}

func RegisterRoute(ar *adminRepo, r *mux.Router) {
	ah := NewAdminHandler(ar)
	r.HandleFunc("/login", ah.Login).Methods("POST")
	return
}

func (handler *APIAdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	thisAdmin := Admin{
		ID:       "",
		Name:     r.FormValue("Name"),
		Password: r.FormValue("Password"),
	}
	storedAdmin, err, message, code := handler.repo.Login(thisAdmin)
	if err != nil {
		api.NewResponse(w, false, err, message, nil, code)
		return
	}
	tkn := jwt.TokenGenerate(w, storedAdmin.ID, storedAdmin.Name)
	if tkn.Authorization == true {
		api.NewResponse(w, true, nil, message, storedAdmin, 200)
	} else {
		api.NewResponse(w, false, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}
