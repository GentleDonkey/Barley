package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	//"gorm.io/driver/mysql"
	//"gorm.io/gorm"
	//_ "gorm.io/gorm"
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
	ID          string `json:"ID"`
	UserID      string `json:"UserID"`
	Description string `json:"Description"` //including purchase date and products, quantity
	Tracking    string `json:"Tracking"`
	Comment     string `json:"Comment"`
	Date        string `json:"Date"` //date of the creation of tracking number
}

type user struct {
	ID         string `json:"ID"`
	WeChatID   string `json:"WeChatID"`
	WeChatName string `json:"WeChatName"`
	//BcryptCode string `json:"BcryptCode"`
}

type admin struct {
	ID       string `json:"ID"`
	Name     string `json:"Name"`
	Password string `json:"Password"`
}

type HttpResponse struct {
	Error   error       `json:"Error"`
	Message string      `json:"Message"`
	Result  interface{} `json:"Result"`
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

	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//db.Raw("SELECT id, name, age FROM users WHERE name = ?", 3).Scan(&result)

	cursor, err := db.Query("SELECT * FROM shipment ORDER BY id DESC")
	if err != nil {
		newResponse(w, err, "Invalid SQL query.", nil, 404)
		return
	}
	var result []shipment
	for cursor.Next() {
		var id, userID, description, tracking, comment, date string
		err = cursor.Scan(&id, &userID, &description, &tracking, &comment, &date)
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
	cursor, err := db.Prepare("INSERT INTO shipment(id, UserID, Description, Tracking, Comment, Date) VALUES(?,?,?,?,?,?)")
	if err != nil {
		newResponse(w, err, "Invalid SQL query.", nil, 404)
		return
	}
	_, err = cursor.Exec(r.FormValue("id"), r.FormValue("UserID"), r.FormValue("Description"), r.FormValue("Tracking"), r.FormValue("Comment"), r.FormValue("Date"))
	if err != nil {
		newResponse(w, err, "Database query error.", nil, 404)
		return
	}
	result := &shipment{
		ID:          r.FormValue("id"),
		UserID:      r.FormValue("UserID"),
		Description: r.FormValue("Description"),
		Tracking:    r.FormValue("Tracking"),
		Comment:     r.FormValue("Comment"),
		Date:        r.FormValue("Date"),
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
	cursor, err := db.Query("SELECT * FROM shipment WHERE id=?", shipmentID)
	if err != nil {
		newResponse(w, err, "Invalid SQL query.", nil, 404)
		return
	}
	var result []shipment
	for cursor.Next() {
		var id, userID, description, tracking, comment, date string
		err = cursor.Scan(&id, &userID, &description, &tracking, &comment, &date)
		if err != nil {
			newResponse(w, err, "Database query error.", nil, 404)
			return
		}
		result = append(result, shipment{ID: id, UserID: userID, Description: description, Tracking: tracking, Comment: comment, Date: date})
	}
	if result == nil {
		err := errors.New("none result error")
		newResponse(w, err, "Database does not found any result.", nil, 404)
		return
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
	cursor, err := db.Prepare("DELETE FROM shipment WHERE id=?")
	if err != nil {
		newResponse(w, err, "Invalid SQL query.", nil, 404)
		return
	}
	_, err = cursor.Exec(shipmentID)
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
	cursor, err := db.Prepare("UPDATE shipment SET UserID=?, Description=?, Tracking=?, Comment=?, Date=? WHERE id=?")
	if err != nil {
		newResponse(w, err, "Invalid SQL query.", nil, 404)
		return
	}
	_, err = cursor.Exec(r.FormValue("UserID"), r.FormValue("Description"), r.FormValue("Tracking"), r.FormValue("Comment"), r.FormValue("Date"), shipmentID)
	if err != nil {
		newResponse(w, err, "Database query error.", nil, 404)
		return
	}
	result := &shipment{
		ID:          shipmentID,
		UserID:      r.FormValue("UserID"),
		Description: r.FormValue("Description"),
		Tracking:    r.FormValue("Tracking"),
		Comment:     r.FormValue("Comment"),
		Date:        r.FormValue("Date"),
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
	cursor, err := db.Query("SELECT * FROM user ORDER BY id")
	if err != nil {
		newResponse(w, err, "Invalid SQL query.", nil, 404)
		return
	}
	var result []user
	for cursor.Next() {
		var id, weChatID, weChatName string
		err = cursor.Scan(&id, &weChatID, &weChatName)
		if err != nil {
			newResponse(w, err, "Database query error.", nil, 404)
			return
		}
		result = append(result, user{ID: id, WeChatID: weChatID, WeChatName: weChatName})
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
	cursor, err := db.Prepare("INSERT INTO user(id, WeChatID, WeChatName) VALUES(?,?,?)")
	if err != nil {
		newResponse(w, err, "Invalid SQL query.", nil, 404)
		return
	}
	_, err = cursor.Exec(r.FormValue("id"), r.FormValue("WeChatID"), r.FormValue("WeChatName"))
	if err != nil {
		newResponse(w, err, "Database query error.", nil, 404)
		return
	}
	result := &user{
		ID:         r.FormValue("id"),
		WeChatID:   r.FormValue("WeChatID"),
		WeChatName: r.FormValue("WeChatName"),
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

/* Admin: Login */
// To login to an admin account
func adminLogin(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	thisAdmin := &admin{
		ID:       "",
		Name:     r.FormValue("Name"),
		Password: r.FormValue("Password"),
	}
	cursor, err := db.Query("SELECT * FROM `admin` WHERE Name=?", thisAdmin.Name)
	if err != nil {
		newResponse(w, err, "Invalid SQL query.", nil, 404)
		return
	}
	storedAdmin := &admin{}
	for cursor.Next() {
		err = cursor.Scan(&storedAdmin.ID, &storedAdmin.Name, &storedAdmin.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				newResponse(w, err, "Admin does not exist.", nil, 401)
				return
			}
			newResponse(w, err, "Database query error.", nil, 404)
			return
		}
	}
	if err = bcrypt.CompareHashAndPassword([]byte(storedAdmin.Password), []byte(thisAdmin.Password)); err != nil {
		newResponse(w, err, "Admin name does not match with password.", nil, 401)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			newResponse(w, err, "Database closing error.", nil, 404)
			return
		}
	}(db)
	newResponse(w, err, "", storedAdmin, 200)
}

/* END-Admin: Login */

/* User: Login */
// To login to a user account
func userLogin(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	userID := mux.Vars(r)["id"]
	_, err := strconv.Atoi(userID)
	if userID == "" || err != nil {
		newResponse(w, err, "Invalid link.", nil, 400)
		return
	}
	thisUser := &user{
		ID:         "",
		WeChatID:   r.FormValue("WeChatID"),
		WeChatName: "",
	}
	cursor, err := db.Query("SELECT * FROM `user` WHERE WeChatID=?", thisUser.WeChatID)
	if err != nil {
		newResponse(w, err, "Invalid SQL query.", nil, 404)
		return
	}
	storedUser := &user{}
	for cursor.Next() {
		err = cursor.Scan(&storedUser.ID, &storedUser.WeChatID, &storedUser.WeChatName)
		if err != nil {
			if err == sql.ErrNoRows {
				newResponse(w, err, "User does not exist.", nil, 401)
				return
			}
			newResponse(w, err, "Database query error.", nil, 404)
			return
		}
	}
	if userID != storedUser.ID {
		newResponse(w, err, "User WeChat ID does not match with user ID.", nil, 401)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			newResponse(w, err, "Database closing error.", nil, 404)
			return
		}
	}(db)
	newResponse(w, err, "", storedUser, 200)
}

/* END-User: Login */

/* User: view shipment */
// To view all shipments
func userAllShipments(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//db.Raw("SELECT id, name, age FROM users WHERE name = ?", 3).Scan(&result)
	userID := mux.Vars(r)["id"]
	_, err := strconv.Atoi(userID)
	if userID == "" || err != nil {
		newResponse(w, err, "Invalid id.", nil, 400)
		return
	}
	cursor, err := db.Query("SELECT * FROM shipment WHERE UserID=?", userID)
	if err != nil {
		newResponse(w, err, "Invalid SQL query.", nil, 404)
		return
	}
	var result []shipment
	for cursor.Next() {
		var id, userID, description, tracking, comment, date string
		err = cursor.Scan(&id, &userID, &description, &tracking, &comment, &date)
		if err != nil {
			newResponse(w, err, "Database query error.", nil, 404)
			return
		}
		result = append(result, shipment{ID: id, UserID: userID, Description: description, Tracking: tracking, Comment: comment, Date: date})
	}
	if result == nil {
		err := errors.New("none result error")
		newResponse(w, err, "Database does not found any result.", nil, 404)
		return
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

/* END-User: view shipment */

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
	router.HandleFunc("/api/v1/admin/login", adminLogin).Methods("POST")
	// User
	router.HandleFunc("/api/v1/user/login/{id}", userLogin).Methods("POST")
	router.HandleFunc("/api/v1/user/{id}", userAllShipments).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
