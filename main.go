package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	//"io/ioutil"
	"log"
	"net/http"
)

//set  database
func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "19900718qzyQZY"
	dbName := "barley"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

type shipment struct {
	ID          string `json:"id"`
	UserID      string `json:"userID"`
	Description string `json:"description"` //including purchase date and products, quantity
	Tracking    string `json:"tracking"`
	Comment     string `json:"comment"`
	Date        string `json:"date"` //date of the creation of tracking number
}

//set router
func firstLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
	//db := dbConn()
	//defer db.Close()
}

func getAllShipments(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM Shipment ORDER BY date DESC")
	if err != nil {
		panic(err.Error())
	}
	var result []shipment
	for selDB.Next() {
		var id, userID, description, tracking, comment, date string
		err = selDB.Scan(&id, &userID, &description, &tracking, &comment, &date)
		if err != nil {
			panic(err.Error())
		}
		result = append(result, shipment{ID: id, UserID: userID, Description: description, Tracking: tracking, Comment: comment, Date: date})
	}
	json.NewEncoder(w).Encode(result)
	defer db.Close()
}

func getOneShipment(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	shipmentID := mux.Vars(r)["id"]
	selDB, err := db.Query("SELECT * FROM Shipment WHERE id=?", shipmentID)
	if err != nil {
		panic(err.Error())
	}
	var result []shipment
	for selDB.Next() {
		var id, userID, description, tracking, comment, date string
		err = selDB.Scan(&id, &userID, &description, &tracking, &comment, &date)
		if err != nil {
			panic(err.Error())
		}
		result = append(result, shipment{ID: id, UserID: userID, Description: description, Tracking: tracking, Comment: comment, Date: date})
	}
	json.NewEncoder(w).Encode(result)
	defer db.Close()
}

func createShipment(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var newShipment shipment
	newShipment.ID = r.FormValue("id")
	newShipment.UserID = r.FormValue("userid")
	newShipment.Description = r.FormValue("description")
	newShipment.Tracking = r.FormValue("tracking")
	newShipment.Comment = r.FormValue("comment")
	newShipment.Date = r.FormValue("date")
	insDB, err := db.Prepare("INSERT INTO shipment(id, userid, description, tracking, comment, date) VALUES(?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	insDB.Exec(newShipment.ID, newShipment.UserID, newShipment.Description, newShipment.Tracking, newShipment.Comment, newShipment.Date)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newShipment)
	defer db.Close()
}

func deleteShipment(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	shipmentID := mux.Vars(r)["id"]
	delDB, err := db.Prepare("DELETE FROM Shipment WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delDB.Exec(shipmentID)
	defer db.Close()
	fmt.Fprintf(w, "The event with ID %v has been deleted successfully", shipmentID)
}

func updateShipment(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	var updatedShipment shipment
	updatedShipment.UserID = r.FormValue("userid")
	updatedShipment.Description = r.FormValue("description")
	updatedShipment.Tracking = r.FormValue("tracking")
	updatedShipment.Comment = r.FormValue("comment")
	updatedShipment.Date = r.FormValue("date")
	shipmentID := mux.Vars(r)["id"]
	updDB, err := db.Prepare("UPDATE Shipment SET userID=?, description=?, tracking=?, comment=?, date=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	updDB.Exec(updatedShipment.UserID, updatedShipment.Description, updatedShipment.Tracking, updatedShipment.Comment, updatedShipment.Date, shipmentID)
	json.NewEncoder(w).Encode(updateShipment)
	defer db.Close()
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", firstLink)
	router.HandleFunc("/admin/shipments", getAllShipments).Methods("GET")
	router.HandleFunc("/admin/shipments/{id}", getOneShipment).Methods("GET")
	router.HandleFunc("/admin/shipment", createShipment).Methods("POST")
	router.HandleFunc("/admin/shipments/{id}", deleteShipment).Methods("DELETE")
	router.HandleFunc("/admin/shipments/{id}", updateShipment).Methods("PATCH")

	log.Fatal(http.ListenAndServe(":8080", router))

}
