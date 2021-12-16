package main

import (
	"database/sql"
	"encoding/json"
	"errors"
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

type user struct {
	ID         string `json:"id"`
	WeChatID   string `json:"wechatID"`
	WeChatName string `json:"wechatname"` //including purchase date and products, quantity
}

type HttpResponse struct {
	Error   error       `json:"error"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

func newResponse(w http.ResponseWriter, error error, message string, result interface{}, statusCode int) {
	var newMessage string
	if error != nil {
		newMessage = strconv.Itoa(statusCode) + ":" + message
	} else {
		newMessage = message
	}
	//fmt.Println(newMessage)
	newResponse := HttpResponse{
		Error:   error,
		Message: newMessage,
		Result:  result,
	}
	jsonNewResp, err := json.Marshal(newResponse)
	if err != nil {
		temp := HttpResponse{
			Error:   errors.New("encode response error"),
			Message: "",
			Result:  nil,
		}
		temp1, _ := json.Marshal(temp)
		w.WriteHeader(500)
		w.Write(temp1)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(jsonNewResp)
}

/* Admin: manage shipment */
// To view all shipments
func getAllShipments(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM shipment ORDER BY id DESC")
	if err != nil {
		newResponse(w, err, "Invalid SQL query.", nil, 404)
		return
	}
	var result []shipment
	for selDB.Next() {
		var id, userID, description, tracking, comment, date string
		err = selDB.Scan(&id, &userID, &description, &tracking, &comment, &date)
		if err != nil {
			newResponse(w, err, "Database query error.", nil, 404)
			return
		}
		result = append(result, shipment{ID: id, UserID: userID, Description: description, Tracking: tracking, Comment: comment, Date: date})
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			newResponse(w, err, "Database closing error.", nil, 404)
			return
		}
	}(db)
	newResponse(w, err, "", result, 200)
}

// To create a new shipment
func createShipment(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	insDB, err := db.Prepare("INSERT INTO shipment(id, userid, description, tracking, comment, date) VALUES(?,?,?,?,?,?)")
	if err != nil {
		newResponse(w, err, "Invalid SQL query.", nil, 404)
		return
	}
	_, err = insDB.Exec(r.FormValue("id"), r.FormValue("userid"), r.FormValue("description"), r.FormValue("tracking"), r.FormValue("comment"), r.FormValue("date"))
	if err != nil {
		newResponse(w, err, "Database query error.", nil, 404)
		return
	}
	result := &shipment{
		ID:          r.FormValue("id"),
		UserID:      r.FormValue("userid"),
		Description: r.FormValue("description"),
		Tracking:    r.FormValue("tracking"),
		Comment:     r.FormValue("comment"),
		Date:        r.FormValue("date"),
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			newResponse(w, err, "Database closing error.", nil, 404)
			return
		}
	}(db)
	newResponse(w, err, "", result, 201)
}

// To view one shipment
func getOneShipment(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	shipmentID := mux.Vars(r)["id"]
	_, err := strconv.Atoi(shipmentID)
	if shipmentID == "" || err != nil {
		newResponse(w, err, "Invalid id.", nil, 400)
		return
	}
	selDB, err := db.Query("SELECT * FROM shipment WHERE id=?", shipmentID)
	if err != nil {
		newResponse(w, err, "Invalid SQL query.", nil, 404)
		return
	}
	var result []shipment
	for selDB.Next() {
		var id, userID, description, tracking, comment, date string
		err = selDB.Scan(&id, &userID, &description, &tracking, &comment, &date)
		if err != nil {
			newResponse(w, err, "Database query error.", nil, 404)
			return
		}
		result = append(result, shipment{ID: id, UserID: userID, Description: description, Tracking: tracking, Comment: comment, Date: date})
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			newResponse(w, err, "Database closing error.", nil, 404)
			return
		}
	}(db)
	newResponse(w, err, "", result, 200)
}

// To delete one shipment
func deleteShipment(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	shipmentID := mux.Vars(r)["id"]
	_, err := strconv.Atoi(shipmentID)
	if shipmentID == "" || err != nil {
		newResponse(w, err, "Invalid id.", nil, 400)
		return
	}
	delDB, err := db.Prepare("DELETE FROM shipment WHERE id=?")
	if err != nil {
		newResponse(w, err, "Invalid SQL query.", nil, 404)
		return
	}
	_, err = delDB.Exec(shipmentID)
	if err != nil {
		newResponse(w, err, "Database query error.", nil, 404)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			newResponse(w, err, "Database closing error.", nil, 404)
			return
		}
	}(db)
	var message = "The shipment with ID " + shipmentID + " has been deleted successfully"
	//fmt.Println(message)
	newResponse(w, err, message, nil, 200)
}

// To update one shipment
func updateShipment(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	shipmentID := mux.Vars(r)["id"]
	_, err := strconv.Atoi(shipmentID)
	if shipmentID == "" || err != nil {
		newResponse(w, err, "Invalid id.", nil, 400)
		return
	}
	updDB, err := db.Prepare("UPDATE shipment SET userID=?, description=?, tracking=?, comment=?, date=? WHERE id=?")
	if err != nil {
		newResponse(w, err, "Invalid SQL query.", nil, 404)
		return
	}
	_, err = updDB.Exec(r.FormValue("userid"), r.FormValue("description"), r.FormValue("tracking"), r.FormValue("comment"), r.FormValue("date"), shipmentID)
	if err != nil {
		newResponse(w, err, "Database query error.", nil, 404)
		return
	}
	result := &shipment{
		ID:          shipmentID,
		UserID:      r.FormValue("userid"),
		Description: r.FormValue("description"),
		Tracking:    r.FormValue("tracking"),
		Comment:     r.FormValue("comment"),
		Date:        r.FormValue("date"),
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			newResponse(w, err, "Database closing error.", nil, 404)
			return
		}
	}(db)
	var message = "The shipment with ID " + shipmentID + " has been updated successfully"
	newResponse(w, err, message, result, 200)
}

/* END-Admin: manage shipment */

/* Admin: manage user account */
// To view all users
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT * FROM user ORDER BY id ASC")
	if err != nil {
		newResponse(w, err, "Invalid SQL query.", nil, 404)
		return
	}
	var result []user
	for selDB.Next() {
		var id, wechatID, wechatname string
		err = selDB.Scan(&id, &wechatID, &wechatname)
		if err != nil {
			newResponse(w, err, "Database query error.", nil, 404)
			return
		}
		result = append(result, user{ID: id, WeChatID: wechatID, WeChatName: wechatname})
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			newResponse(w, err, "Database closing error.", nil, 404)
			return
		}
	}(db)
	newResponse(w, err, "", result, 200)
}

// To create a new user account
func createUser(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	insDB, err := db.Prepare("INSERT INTO user(id, WeChatID, WeChatName) VALUES(?,?,?)")
	if err != nil {
		newResponse(w, err, "Invalid SQL query.", nil, 404)
		return
	}
	_, err = insDB.Exec(r.FormValue("id"), r.FormValue("wechatid"), r.FormValue("wechatname"))
	if err != nil {
		newResponse(w, err, "Database query error.", nil, 404)
		return
	}
	result := &user{
		ID:         r.FormValue("id"),
		WeChatID:   r.FormValue("wechatid"),
		WeChatName: r.FormValue("wechatname"),
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			newResponse(w, err, "Database closing error.", nil, 404)
			return
		}
	}(db)
	newResponse(w, err, "", result, 201)
}

/* END-Admin: manage user account */

func main() {
	router := mux.NewRouter().StrictSlash(true)
	// Admin
	router.HandleFunc("/api/v1/admin/shipments", getAllShipments).Methods("GET")
	router.HandleFunc("/api/v1/admin/shipment", createShipment).Methods("POST")
	router.HandleFunc("/api/v1/admin/shipments/{id}", getOneShipment).Methods("GET")
	router.HandleFunc("/api/v1/admin/shipments/{id}", deleteShipment).Methods("DELETE")
	router.HandleFunc("/api/v1/admin/shipments/{id}", updateShipment).Methods("PATCH")
	router.HandleFunc("/api/v1/admin/users", getAllUsers).Methods("GET")
	router.HandleFunc("/api/v1/admin/user", createUser).Methods("POST")
	//User
	log.Fatal(http.ListenAndServe(":8080", router))
}
