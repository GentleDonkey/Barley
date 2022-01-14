package main

import (
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"notifications/configs"
	"notifications/pkg/server"
)

func main() {
	r := server.SetServer()
	http.ListenAndServe(configs.Host, r)
}
