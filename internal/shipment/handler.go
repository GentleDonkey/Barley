package shipment

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"notifications/internal/api"
	myError "notifications/internal/error"
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
	var shipment Shipment
	err := json.NewDecoder(r.Body).Decode(&shipment)
	if err != nil {
		api.NewHttpResponse(w, myError.NewError(err, "Unable to convert a json to an object", 500), "", nil)
		return
	}
	repoErr := handler.repo.Create(shipment)
	if repoErr != nil {
		api.NewHttpResponse(w, repoErr, "", nil)
		return
	}
	newMessage := "201: A new shipment with ID " + shipment.ID + " has been created successfully"
	api.NewHttpResponse(w, nil, newMessage, nil)
	return
}

func (handler *APIShipmentHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	result, repoErr := handler.repo.FindAll()
	if repoErr != nil {
		api.NewHttpResponse(w, repoErr, "", nil)
		return
	}
	newMessage := "200: All shipment have been found successfully"
	api.NewHttpResponse(w, nil, newMessage, result)
	return
}

func (handler *APIShipmentHandler) FindOne(w http.ResponseWriter, r *http.Request) {
	shipmentID := mux.Vars(r)["id"]
	_, err := strconv.Atoi(shipmentID)
	if shipmentID == "" || err != nil {
		api.NewHttpResponse(w, &myError.InvalidPara, "", nil)
		return
	}
	result := handler.repo.FindOne(shipmentID)
	newMessage := "200: A new shipment with ID " + shipmentID + " has been found successfully"
	api.NewHttpResponse(w, nil, newMessage, result)
	return
}

func (handler *APIShipmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	shipmentID := mux.Vars(r)["id"]
	_, err := strconv.Atoi(shipmentID)
	if shipmentID == "" || err != nil {
		api.NewHttpResponse(w, &myError.InvalidPara, "", nil)
		return
	}
	var shipment Shipment
	err = json.NewDecoder(r.Body).Decode(&shipment)
	if err != nil {
		api.NewHttpResponse(w, myError.NewError(err, "Unable to convert a json to an object", 500), "", nil)
		return
	}
	shipment.ID = shipmentID
	repoErr := handler.repo.Update(shipment)
	if repoErr != nil {
		api.NewHttpResponse(w, repoErr, "", nil)
		return
	}
	newMessage := "200: The shipment with ID " + shipment.ID + " has been updated successfully"
	api.NewHttpResponse(w, nil, newMessage, nil)
	return
}

func (handler *APIShipmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	shipmentID := mux.Vars(r)["id"]
	_, err := strconv.Atoi(shipmentID)
	if shipmentID == "" || err != nil {
		api.NewHttpResponse(w, &myError.InvalidPara, "", nil)
		return
	}
	repoErr := handler.repo.Delete(shipmentID)
	if repoErr != nil {
		api.NewHttpResponse(w, repoErr, "", nil)
		return
	}
	newMessage := "200: The shipment with ID " + shipmentID + " has been deleted successfully"
	api.NewHttpResponse(w, nil, newMessage, nil)
}
