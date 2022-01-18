package shipment

import (
	"github.com/gorilla/mux"
	"net/http"
	"notifications/pkg/api"
	"notifications/pkg/jwt"
	"strconv"
)

type APIShipmentHandler struct {
	repo *shipmentRepo
}

func NewShipmentHandler(sr *shipmentRepo) *APIShipmentHandler {
	return &APIShipmentHandler{
		sr,
	}
}

func RegisterRoute(sr *shipmentRepo, r *mux.Router) {
	sh := NewShipmentHandler(sr)
	r.HandleFunc("/shipment", sh.Create).Methods("POST")
	r.HandleFunc("/shipments", sh.FindAll).Methods("GET")
	r.HandleFunc("/shipment/{id}/", sh.FindOne).Methods("GET")
	r.HandleFunc("/shipment/{id}", sh.Delete).Methods("DELETE")
	r.HandleFunc("/shipment/{id}", sh.Update).Methods("PATCH")
}

func (handler *APIShipmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	tkn := jwt.TokenParse(r)
	if tkn.Authorization == true {
		shipment := Shipment{
			ID:          r.FormValue("ID"),
			UserID:      r.FormValue("UserID"),
			Description: r.FormValue("Description"),
			Tracking:    r.FormValue("Tracking"),
			Comment:     r.FormValue("Comment"),
			Date:        r.FormValue("Date"),
		}
		err, message, code := handler.repo.Create(shipment)
		api.NewResponse(w, true, err, message, shipment, code)
	} else {
		api.NewResponse(w, false, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}

func (handler *APIShipmentHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	tkn := jwt.TokenParse(r)
	if tkn.Authorization == true {
		result, err, message, code := handler.repo.FindAll()
		api.NewResponse(w, true, err, message, result, code)
	} else {
		api.NewResponse(w, false, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}

func (handler *APIShipmentHandler) FindOne(w http.ResponseWriter, r *http.Request) {
	tkn := jwt.TokenParse(r)
	if tkn.Authorization == true {
		shipmentID := mux.Vars(r)["id"]
		_, err := strconv.Atoi(shipmentID)
		if shipmentID == "" || err != nil {
			api.NewResponse(w, tkn.Authorization, err, "Invalid id.", nil, 400)
			return
		}
		result, err, message, code := handler.repo.FindOne(shipmentID)
		if err != nil {
			api.NewResponse(w, true, err, message, nil, code)
			return
		}
		api.NewResponse(w, true, nil, message, result, code)
	} else {
		api.NewResponse(w, tkn.Authorization, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}

func (handler *APIShipmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	tkn := jwt.TokenParse(r)
	if tkn.Authorization == true {
		shipmentID := mux.Vars(r)["id"]
		_, err := strconv.Atoi(shipmentID)
		if shipmentID == "" || err != nil {
			api.NewResponse(w, tkn.Authorization, err, "Invalid id.", nil, 400)
			return
		}
		shipment := Shipment{
			ID:          shipmentID,
			UserID:      r.FormValue("UserID"),
			Description: r.FormValue("Description"),
			Tracking:    r.FormValue("Tracking"),
			Comment:     r.FormValue("Comment"),
			Date:        r.FormValue("Date"),
		}
		err, message, code := handler.repo.Update(shipment)
		api.NewResponse(w, true, err, message, shipment, code)
	} else {
		api.NewResponse(w, tkn.Authorization, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}

func (handler *APIShipmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	tkn := jwt.TokenParse(r)
	if tkn.Authorization == true {
		shipmentID := mux.Vars(r)["id"]
		_, err := strconv.Atoi(shipmentID)
		if shipmentID == "" || err != nil {
			api.NewResponse(w, tkn.Authorization, err, "Invalid id.", nil, 400)
			return
		}
		err, message, code := handler.repo.Delete(shipmentID)
		api.NewResponse(w, true, err, message, nil, code)
	} else {
		api.NewResponse(w, tkn.Authorization, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}
