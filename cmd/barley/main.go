package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"log"
	"net/http"
	"notifications/configs"
	"notifications/internal/server"
)

func main() {
	config := configs.LoadConfig()
	r := server.SetServer(config)
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Origin"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe(config.ServerAddress, handlers.CORS(originsOk, headersOk, methodsOk)(r)))
}
