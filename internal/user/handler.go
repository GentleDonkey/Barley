package user

import (
	"github.com/gorilla/mux"
	"net/http"
	"notifications/pkg/api"
	"notifications/pkg/jwt"
)

type APIUserHandler struct {
	repo *userRepo
}

func NewUserHandler(ur *userRepo) *APIUserHandler {
	return &APIUserHandler{
		ur,
	}
}

func RegisterRoute(ur *userRepo, r *mux.Router) {
	uh := NewUserHandler(ur)
	r.HandleFunc("/user", uh.Create).Methods("POST")
	r.HandleFunc("/users", uh.FindAll).Methods("GET")
}

func (handler *APIUserHandler) Create(w http.ResponseWriter, r *http.Request) {
	tkn := jwt.TokenParse(r)
	if tkn.Authorization == true {
		randomCode := api.RandomString(16)
		user := User{
			ID:         r.FormValue("ID"),
			WeChatID:   r.FormValue("WeChatID"),
			WeChatName: r.FormValue("WeChatName"),
			RandomCode: randomCode,
		}
		err, message, code := handler.repo.Create(user)
		api.NewResponse(w, true, err, message, user, code)
	} else {
		api.NewResponse(w, false, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}

func (handler *APIUserHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	tkn := jwt.TokenParse(r)
	if tkn.Authorization == true {
		result, err, message, code := handler.repo.FindAll()
		api.NewResponse(w, true, err, message, result, code)
	} else {
		api.NewResponse(w, false, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}
