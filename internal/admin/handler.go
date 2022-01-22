package admin

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"notifications/internal/api"
	myError "notifications/internal/error"
	"notifications/internal/jwt"
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
	var thisAdmin Admin
	err := json.NewDecoder(r.Body).Decode(&thisAdmin)
	if err != nil {
		api.NewHttpResponse(w, myError.NewError(err, "Unable to convert a json to an object", 500), "", nil)
		return
	}
	storedAdmin, repoErr := handler.repo.Login(thisAdmin)
	if repoErr != nil {
		api.NewHttpResponse(w, repoErr, "", nil)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(storedAdmin.Password), []byte(thisAdmin.Password))
	if err != nil {
		api.NewHttpResponse(w, myError.NewError(err, "Admin name does not match with password.", 401), "", nil)
		return
	}
	if jwt.TokenGenerate(w, storedAdmin.ID, storedAdmin.Name) == false {
		api.NewHttpResponse(w, myError.NewError(errors.New("JTW error"), "Failed to generate a JWT token.", 401), "", nil)
		return
	}
	message := "200: Admin " + thisAdmin.Name + " login successfully"
	api.NewHttpResponse(w, nil, message, nil)
	return
}
