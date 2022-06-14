package tracking

import (
	"github.com/gorilla/mux"
	"net/http"
	"notifications/internal/api"
	myError "notifications/internal/error"
)

type APITrackingHandler struct {
	repo APITrackingRepo
}

func NewTrackingHandler(sr APITrackingRepo) *APITrackingHandler {
	return &APITrackingHandler{
		sr,
	}
}

func RegisterRoute(sr APITrackingRepo, r *mux.Router) {
	th := NewTrackingHandler(sr)
	r.HandleFunc("/{code}/shipments", th.FindAll).Methods("GET")
}

func (handler *APITrackingHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	randomCode := mux.Vars(r)["code"]
	if randomCode == "" {
		api.NewHttpResponse(w, &myError.InvalidPara, "", nil)
		return
	}
	result, repoErr := handler.repo.FindAll(randomCode)
	if repoErr != nil {
		api.NewHttpResponse(w, repoErr, "", nil)
		return
	}
	newMessage := "200: All shipment have been found successfully"
	api.NewHttpResponse(w, nil, newMessage, result)
	return
}
