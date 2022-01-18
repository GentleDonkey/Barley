package tracking

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"notifications/pkg/api"
)

type APITrackingHandler struct {
	repo *trackingRepo
}

func NewTrackingHandler(sr *trackingRepo) *APITrackingHandler {
	return &APITrackingHandler{
		sr,
	}
}

func RegisterRoute(sr *trackingRepo, r *mux.Router) {
	th := NewTrackingHandler(sr)
	r.HandleFunc("/tracking/{code}", th.FindAll).Methods("GET")
}

func (handler *APITrackingHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	randomCode := mux.Vars(r)["code"]
	if randomCode == "" {
		api.NewResponse(w, false, errors.New("empty code"), "Empty code parameter.", nil, 400)
		return
	}
	result, err, message, code := handler.repo.FindAll(randomCode)
	api.NewResponse(w, true, err, message, result, code)
}
