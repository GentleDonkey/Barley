package user

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"notifications/internal/api"
	myError "notifications/internal/error"
	"time"
)

type APIUserHandler struct {
	repo *userRepo
}

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

func NewUserHandler(ur *userRepo) *APIUserHandler {
	return &APIUserHandler{
		ur,
	}
}

func RegisterRoute(ur *userRepo, r *mux.Router) {
	uh := NewUserHandler(ur)
	r.HandleFunc("/user", uh.Create).Methods("POST")
	r.HandleFunc("/users", uh.FindAll).Methods("GET")
}

func (handler *APIUserHandler) Create(w http.ResponseWriter, r *http.Request) {
	randomCode := RandomString(16)
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		api.NewHttpResponse(w, myError.NewError(err, "Unable to convert a json to an object", 500), "", nil)
		return
	}
	user.RandomCode = randomCode
	repoErr := handler.repo.Create(user)
	if repoErr != nil {
		api.NewHttpResponse(w, repoErr, "", nil)
		return
	}
	newMessage := "201: A new user with ID " + user.ID + " has been created successfully"
	api.NewHttpResponse(w, nil, newMessage, nil)
	return
}

func (handler *APIUserHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	result, repoErr := handler.repo.FindAll()
	if repoErr != nil {
		api.NewHttpResponse(w, repoErr, "", nil)
		return
	}
	newMessage := "200: All user have been found successfully"
	api.NewHttpResponse(w, nil, newMessage, result)
	return
}
