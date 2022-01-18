package main

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"notifications/configs"
	"notifications/pkg/server"
)

func main() {
	r := server.SetServer()

	log.Fatal(http.ListenAndServe(configs.Host+configs.RestPath, r))
}
