package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
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
	RandomCode string `json:"RandomCode"`
}

type admin struct {
	ID       string `json:"ID"`
	Name     string `json:"Name"`
	Password string `json:"Password"`
}

type HttpResponse struct {
	Authorization bool        `json:"Authorization"`
	Error         error       `json:"Error"`
	Message       string      `json:"Message"`
	Result        interface{} `json:"Result"`
	StatusCode    int         `json:"StatusCode"`
}

type Claims struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	jwt.StandardClaims
}

var jwtKey = []byte("my_secret_key")

func RandomString(length int) string {
	var seededRand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func newResponse(w http.ResponseWriter, authorization bool, error error, message string, result interface{}, statusCode int) {
	var newMessage string
	if error != nil {
		newMessage = strconv.Itoa(statusCode) + ":" + message
	} else {
		newMessage = message
	}
	newResponse := HttpResponse{
		Authorization: authorization,
		Error:         error,
		Message:       newMessage,
		Result:        result,
		StatusCode:    statusCode,
	}
	jsonNewResp, err := json.Marshal(newResponse)
	if err != nil {
		temp := HttpResponse{
			Authorization: authorization,
			Error:         errors.New("encode response error"),
			Message:       "",
			Result:        nil,
			StatusCode:    500,
		}
		temp1, _ := json.Marshal(temp)
		w.WriteHeader(500)
		w.Write(temp1)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(jsonNewResp)
}

func tokenGenerate(w http.ResponseWriter, credentialID string, credentialName string) HttpResponse {
	newResponse := HttpResponse{}
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		ID:   credentialID,
		Name: credentialName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		newResponse.Authorization = false
		newResponse.Error = err
		newResponse.Message = "Failed to get the signed JWT."
		newResponse.Result = nil
		newResponse.StatusCode = 500
		return newResponse
	}
	w.Header().Add("Authorization", tokenString)
	newResponse.Authorization = true
	newResponse.Error = nil
	newResponse.Message = ""
	newResponse.Result = tokenString
	newResponse.StatusCode = 200
	return newResponse
}

func tokenParse(r *http.Request) HttpResponse {
	newResponse := HttpResponse{}
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		newResponse.Authorization = false
		newResponse.Error = errors.New("authorization not set")
		newResponse.Message = "Unauthorized: missing a token."
		newResponse.Result = nil
		newResponse.StatusCode = 401
		return newResponse
	}
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			newResponse.Authorization = false
			newResponse.Error = err
			newResponse.Message = "Unauthorized: signature is invalid."
			newResponse.Result = nil
			newResponse.StatusCode = 401
			return newResponse
		}
		newResponse.Authorization = false
		newResponse.Error = err
		newResponse.Message = "Unauthorized: internal error."
		newResponse.Result = nil
		newResponse.StatusCode = 400
		return newResponse
	}
	if !token.Valid {
		newResponse.Authorization = false
		newResponse.Error = err
		newResponse.Message = "Unauthorized: token is invalid."
		newResponse.Result = nil
		newResponse.StatusCode = 400
		return newResponse
	}
	newResponse.Authorization = true
	newResponse.Error = nil
	newResponse.Message = ""
	newResponse.Result = token
	newResponse.StatusCode = 200
	return newResponse
}

/* Admin: manage shipment */
// To view all shipments
func getAllShipments(w http.ResponseWriter, r *http.Request) {
	tkn := tokenParse(r)
	if tkn.Authorization == true {
		db := dbConn()
		cursor, err := db.Query("SELECT * FROM shipment ORDER BY id DESC")
		if err != nil {
			newResponse(w, tkn.Authorization, err, "Invalid SQL query.", nil, 404)
			return
		}
		var result []shipment
		for cursor.Next() {
			var id, userID, description, tracking, comment, date string
			err = cursor.Scan(&id, &userID, &description, &tracking, &comment, &date)
			if err != nil {
				newResponse(w, tkn.Authorization, err, "Database query error.", nil, 404)
				return
			}
			result = append(result, shipment{ID: id, UserID: userID, Description: description, Tracking: tracking, Comment: comment, Date: date})
		}
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				newResponse(w, tkn.Authorization, err, "Database closing error.", nil, 404)
				return
			}
		}(db)
		newResponse(w, tkn.Authorization, err, "", result, 200)
	} else {
		newResponse(w, tkn.Authorization, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}

// To create a new shipment
func createShipment(w http.ResponseWriter, r *http.Request) {
	tkn := tokenParse(r)
	if tkn.Authorization == true {
		db := dbConn()
		cursor, err := db.Prepare("INSERT INTO shipment(id, UserID, Description, Tracking, Comment, Date) VALUES(?,?,?,?,?,?)")
		if err != nil {
			newResponse(w, tkn.Authorization, err, "Invalid SQL query.", nil, 404)
			return
		}
		_, err = cursor.Exec(r.FormValue("id"), r.FormValue("UserID"), r.FormValue("Description"), r.FormValue("Tracking"), r.FormValue("Comment"), r.FormValue("Date"))
		if err != nil {
			newResponse(w, tkn.Authorization, err, "Database query error.", nil, 404)
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
				newResponse(w, tkn.Authorization, err, "Database closing error.", nil, 404)
				return
			}
		}(db)
		newResponse(w, tkn.Authorization, err, "", result, 201)
	} else {
		newResponse(w, tkn.Authorization, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}

// To view one shipment
func getOneShipment(w http.ResponseWriter, r *http.Request) {
	tkn := tokenParse(r)
	if tkn.Authorization == true {
		db := dbConn()
		shipmentID := mux.Vars(r)["id"]
		_, err := strconv.Atoi(shipmentID)
		if shipmentID == "" || err != nil {
			newResponse(w, tkn.Authorization, err, "Invalid id.", nil, 400)
			return
		}
		cursor, err := db.Query("SELECT * FROM shipment WHERE id=?", shipmentID)
		if err != nil {
			newResponse(w, tkn.Authorization, err, "Invalid SQL query.", nil, 404)
			return
		}
		var result []shipment
		for cursor.Next() {
			var id, userID, description, tracking, comment, date string
			err = cursor.Scan(&id, &userID, &description, &tracking, &comment, &date)
			if err != nil {
				newResponse(w, tkn.Authorization, err, "Database query error.", nil, 404)
				return
			}
			result = append(result, shipment{ID: id, UserID: userID, Description: description, Tracking: tracking, Comment: comment, Date: date})
		}
		if result == nil {
			err := errors.New("none result error")
			newResponse(w, tkn.Authorization, err, "Database does not found any result.", nil, 404)
			return
		}
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				newResponse(w, tkn.Authorization, err, "Database closing error.", nil, 404)
				return
			}
		}(db)
		newResponse(w, tkn.Authorization, err, "", result, 200)
	} else {
		newResponse(w, tkn.Authorization, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}

// To delete one shipment
func deleteShipment(w http.ResponseWriter, r *http.Request) {
	tkn := tokenParse(r)
	if tkn.Authorization == true {
		db := dbConn()
		shipmentID := mux.Vars(r)["id"]
		_, err := strconv.Atoi(shipmentID)
		if shipmentID == "" || err != nil {
			newResponse(w, tkn.Authorization, err, "Invalid id.", nil, 400)
			return
		}
		cursor, err := db.Prepare("DELETE FROM shipment WHERE id=?")
		if err != nil {
			newResponse(w, tkn.Authorization, err, "Invalid SQL query.", nil, 404)
			return
		}
		_, err = cursor.Exec(shipmentID)
		if err != nil {
			newResponse(w, tkn.Authorization, err, "Database query error.", nil, 404)
			return
		}
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				newResponse(w, tkn.Authorization, err, "Database closing error.", nil, 404)
				return
			}
		}(db)
		var message = "The shipment with ID " + shipmentID + " has been deleted successfully"
		//fmt.Println(message)
		newResponse(w, tkn.Authorization, err, message, nil, 200)
	} else {
		newResponse(w, tkn.Authorization, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}

// To update one shipment
func updateShipment(w http.ResponseWriter, r *http.Request) {
	tkn := tokenParse(r)
	if tkn.Authorization == true {
		db := dbConn()
		shipmentID := mux.Vars(r)["id"]
		_, err := strconv.Atoi(shipmentID)
		if shipmentID == "" || err != nil {
			newResponse(w, tkn.Authorization, err, "Invalid id.", nil, 400)
			return
		}
		cursor, err := db.Prepare("UPDATE shipment SET UserID=?, Description=?, Tracking=?, Comment=?, Date=? WHERE id=?")
		if err != nil {
			newResponse(w, tkn.Authorization, err, "Invalid SQL query.", nil, 404)
			return
		}
		_, err = cursor.Exec(r.FormValue("UserID"), r.FormValue("Description"), r.FormValue("Tracking"), r.FormValue("Comment"), r.FormValue("Date"), shipmentID)
		if err != nil {
			newResponse(w, tkn.Authorization, err, "Database query error.", nil, 404)
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
				newResponse(w, tkn.Authorization, err, "Database closing error.", nil, 404)
				return
			}
		}(db)
		var message = "The shipment with ID " + shipmentID + " has been updated successfully"
		newResponse(w, tkn.Authorization, err, message, result, 200)
	} else {
		newResponse(w, tkn.Authorization, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}

/* END-Admin: manage shipment */

/* Admin: manage user account */
// To view all users
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	tkn := tokenParse(r)
	if tkn.Authorization == true {
		db := dbConn()
		cursor, err := db.Query("SELECT * FROM user ORDER BY id")
		if err != nil {
			newResponse(w, tkn.Authorization, err, "Invalid SQL query.", nil, 404)
			return
		}
		var result []user
		for cursor.Next() {
			var id, weChatID, weChatName string
			err = cursor.Scan(&id, &weChatID, &weChatName)
			if err != nil {
				newResponse(w, tkn.Authorization, err, "Database query error.", nil, 404)
				return
			}
			result = append(result, user{ID: id, WeChatID: weChatID, WeChatName: weChatName})
		}
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				newResponse(w, tkn.Authorization, err, "Database closing error.", nil, 404)
				return
			}
		}(db)
		newResponse(w, tkn.Authorization, err, "", result, 200)
	} else {
		newResponse(w, tkn.Authorization, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}

// To create a new user account
func createUser(w http.ResponseWriter, r *http.Request) {
	tkn := tokenParse(r)
	if tkn.Authorization == true {
		db := dbConn()
		code := RandomString(16)
		cursor, err := db.Prepare("INSERT INTO user(id, WeChatID, WeChatName, RandomCode) VALUES(?,?,?,?)")
		if err != nil {
			newResponse(w, tkn.Authorization, err, "Invalid SQL query.", nil, 404)
			return
		}
		_, err = cursor.Exec(r.FormValue("id"), r.FormValue("WeChatID"), r.FormValue("WeChatName"), code)
		if err != nil {
			newResponse(w, tkn.Authorization, err, "Database query error.", nil, 404)
			return
		}
		result := &user{
			ID:         r.FormValue("id"),
			WeChatID:   r.FormValue("WeChatID"),
			WeChatName: r.FormValue("WeChatName"),
			RandomCode: code,
		}
		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
				newResponse(w, tkn.Authorization, err, "Database closing error.", nil, 404)
				return
			}
		}(db)
		newResponse(w, tkn.Authorization, err, "", result, 201)
	} else {
		newResponse(w, tkn.Authorization, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
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
		newResponse(w, false, err, "Invalid SQL query.", nil, 404)
		return
	}
	storedAdmin := &admin{}
	for cursor.Next() {
		err = cursor.Scan(&storedAdmin.ID, &storedAdmin.Name, &storedAdmin.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				newResponse(w, false, err, "Admin does not exist.", nil, 401)
				return
			}
			newResponse(w, false, err, "Database query error.", nil, 404)
			return
		}
	}
	if err = bcrypt.CompareHashAndPassword([]byte(storedAdmin.Password), []byte(thisAdmin.Password)); err != nil {
		newResponse(w, false, err, "Admin name does not match with password.", nil, 401)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			newResponse(w, false, err, "Database closing error.", nil, 404)
			return
		}
	}(db)
	tkn := tokenGenerate(w, storedAdmin.ID, storedAdmin.Name)
	if tkn.Authorization == true {
		newResponse(w, true, err, "", storedAdmin, 200)
	} else {
		newResponse(w, tkn.Authorization, tkn.Error, tkn.Message, nil, tkn.StatusCode)
	}
}

/* END-Admin: Login */

/* User: Login */

/* User: view shipments */
// To view all shipments
func userAllShipments(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	code := mux.Vars(r)["code"]
	if code == "" {
		newResponse(w, false, errors.New("empty code"), "Empty code parameter.", nil, 400)
		return
	}
	cursor, err := db.Query("SELECT shipment.* FROM shipment JOIN user ON (user.id=shipment.userID AND user.RandomCode=?)", code)
	if err != nil {
		newResponse(w, false, err, "Invalid SQL query.", nil, 404)
		return
	}
	var result []shipment
	for cursor.Next() {
		var id, userID, description, tracking, comment, date string
		err = cursor.Scan(&id, &userID, &description, &tracking, &comment, &date)
		if err != nil {
			newResponse(w, false, err, "Database query error.", nil, 404)
			return
		}
		result = append(result, shipment{ID: id, UserID: userID, Description: description, Tracking: tracking, Comment: comment, Date: date})
	}
	if result == nil {
		err := errors.New("none result error")
		newResponse(w, false, err, "Database does not found any result.", nil, 404)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			newResponse(w, false, err, "Database closing error.", nil, 404)
			return
		}
	}(db)
	newResponse(w, false, err, "", result, 200)

}

/* END-User: view shipments */

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
	router.HandleFunc("/api/v1/user/tracking/{code}", userAllShipments).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
