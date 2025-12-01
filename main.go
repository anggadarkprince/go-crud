package main

import (
	"log"
	"net/http"
	"os"

	"github.com/anggadarkprince/crud-employee-go/database"
	"github.com/anggadarkprince/crud-employee-go/routes"
	"github.com/anggadarkprince/crud-employee-go/utilities"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}
	
	utilities.InitTemplates()

	db := database.InitDatabase()
	
	server := http.NewServeMux()

	server.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/favicon.ico")
	})
	routes.MapRoutes(server, db)

	app_port := os.Getenv("APP_PORT")
	http.ListenAndServe(":" + app_port, server)
}
