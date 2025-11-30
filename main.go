package main

import (
	"net/http"

	"github.com/anggadarkprince/crud-employee-go/database"
	"github.com/anggadarkprince/crud-employee-go/routes"
	"github.com/anggadarkprince/crud-employee-go/utilities"
)

func main() {
	utilities.InitTemplates()

	db := database.InitDatabase()
	
	server := http.NewServeMux()

	routes.MapRoutes(server, db)

	http.ListenAndServe(":8080", server)
}
