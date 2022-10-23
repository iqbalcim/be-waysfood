package main

import (
	"fmt"
	"net/http"
	"waysfood/database"
	"waysfood/pkg/mysql"
	"waysfood/routes"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	mysql.DatabaseInit()
	database.RunMigration()
	routes.RouteInit(r.PathPrefix("/api/v1").Subrouter())

	fmt.Println("server running localhost:8080")
	http.ListenAndServe("localhost:8080", r)
}