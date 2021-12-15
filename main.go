package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"strconv"

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
		fmt.Println(err.Error())
		return
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
func getAllShipments(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM shipment ORDER BY id DESC")
	if err != nil {
		w.WriteHeader(404)
		fmt.Println(err.Error())
		return
	}
	var result []shipment
	for selDB.Next() {
		var id, userID, description, tracking, comment, date string
		err = selDB.Scan(&id, &userID, &description, &tracking, &comment, &date)
		if err != nil {
			w.WriteHeader(404)
			fmt.Println(err.Error())
			return
		}
		result = append(result, shipment{ID: id, UserID: userID, Description: description, Tracking: tracking, Comment: comment, Date: date})
	}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		w.WriteHeader(404)
		fmt.Println(err.Error())
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			w.WriteHeader(404)
			fmt.Println(err.Error())
			return
		}
	}(db)
	w.WriteHeader(200)
}

func getOneShipment(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	shipmentID := mux.Vars(r)["id"]
	_, err := strconv.Atoi(shipmentID)
	if shipmentID == "" || err != nil {
		w.WriteHeader(401)
		fmt.Println(err.Error())
		return
	}
	selDB, err := db.Query("SELECT * FROM shipment WHERE id=?", shipmentID)
	if err != nil {
		w.WriteHeader(404)
		fmt.Println(err.Error())
		return
	}
	var result []shipment
	for selDB.Next() {
		var id, userID, description, tracking, comment, date string
		err = selDB.Scan(&id, &userID, &description, &tracking, &comment, &date)
		if err != nil {
			w.WriteHeader(404)
			fmt.Println(err.Error())
			return
		}
		result = append(result, shipment{ID: id, UserID: userID, Description: description, Tracking: tracking, Comment: comment, Date: date})
	}
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		w.WriteHeader(404)
		fmt.Println(err.Error())
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			w.WriteHeader(404)
			fmt.Println(err.Error())
			return
		}
	}(db)
	w.WriteHeader(200)
}

func createShipment(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	insDB, err := db.Prepare("INSERT INTO shipment(id, userid, description, tracking, comment, date) VALUES(?,?,?,?,?,?)")
	if err != nil {
		err := fmt.Errorf("name %v", err)
		log.Printf("name %v", err)
		fmt.Println(err.Error())
	}
	_, err = insDB.Exec(r.FormValue("id"), r.FormValue("userid"), r.FormValue("description"), r.FormValue("tracking"), r.FormValue("comment"), r.FormValue("date"))
	if err != nil {
		w.WriteHeader(404)
		fmt.Println(err.Error())
	}
	/* Only for test */
	newShipment := &shipment{
		ID:          r.FormValue("id"),
		UserID:      r.FormValue("userid"),
		Description: r.FormValue("description"),
		Tracking:    r.FormValue("tracking"),
		Comment:     r.FormValue("comment"),
		Date:        r.FormValue("date"),
	}
	/* End */
	err = json.NewEncoder(w).Encode(newShipment)
	if err != nil {
		w.WriteHeader(404)
		fmt.Println(err.Error())
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			w.WriteHeader(404)
			fmt.Println(err.Error())
		}
	}(db)
	w.WriteHeader(201)
}

func deleteShipment(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	shipmentID := mux.Vars(r)["id"]
	_, err := strconv.Atoi(shipmentID)
	if shipmentID == "" || err != nil {
		w.WriteHeader(401)
		fmt.Println(err.Error())
		return
	}
	delDB, err := db.Prepare("DELETE FROM shipment WHERE id=?")
	if err != nil {
		w.WriteHeader(404)
		fmt.Println(err.Error())
		return
	}
	_, err = delDB.Exec(shipmentID)
	if err != nil {
		w.WriteHeader(404)
		fmt.Println(err.Error())
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			w.WriteHeader(404)
			fmt.Println(err.Error())
			return
		}
	}(db)
	fmt.Fprintf(w, "The event with ID %v has been deleted successfully", shipmentID)
	w.WriteHeader(200)
}

func updateShipment(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	shipmentID := mux.Vars(r)["id"]
	_, err := strconv.Atoi(shipmentID)
	if shipmentID == "" || err != nil {
		w.WriteHeader(401)
		fmt.Println(err.Error())
		return
	}
	updDB, err := db.Prepare("UPDATE shipment SET userID=?, description=?, tracking=?, comment=?, date=? WHERE id=?")
	if err != nil {
		w.WriteHeader(404)
		fmt.Println(err.Error())
		return
	}
	_, err = updDB.Exec(shipmentID, r.FormValue("userid"), r.FormValue("description"), r.FormValue("tracking"), r.FormValue("comment"), r.FormValue("date"))
	if err != nil {
		w.WriteHeader(404)
		fmt.Println(err.Error())
		return
	}
	/* Only for test */
	updateShipment := &shipment{
		ID:          shipmentID,
		UserID:      r.FormValue("userid"),
		Description: r.FormValue("description"),
		Tracking:    r.FormValue("tracking"),
		Comment:     r.FormValue("comment"),
		Date:        r.FormValue("date"),
	}
	/* End */
	err = json.NewEncoder(w).Encode(updateShipment)
	if err != nil {
		w.WriteHeader(404)
		fmt.Println(err.Error())
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			w.WriteHeader(404)
			fmt.Println(err.Error())
			return
		}
	}(db)
	fmt.Fprintf(w, "The event with ID %v has been updated successfully", shipmentID)
	w.WriteHeader(200)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/admin/shipments", getAllShipments).Methods("GET")
	router.HandleFunc("/admin/shipments/{id}", getOneShipment).Methods("GET")
	router.HandleFunc("/admin/shipment", createShipment).Methods("POST")
	router.HandleFunc("/admin/shipments/{id}", deleteShipment).Methods("DELETE")
	router.HandleFunc("/admin/shipments/{id}", updateShipment).Methods("PATCH")
	log.Fatal(http.ListenAndServe(":8080", router))
}
