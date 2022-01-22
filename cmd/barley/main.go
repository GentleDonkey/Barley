package main

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"notifications/configs"
	"notifications/internal/server"
)

func main() {
	config := configs.LoadConfig()
	r := server.SetServer(config)
	log.Fatal(http.ListenAndServe(config.ServerAddress, r))
}
