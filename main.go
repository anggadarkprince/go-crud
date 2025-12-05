package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/anggadarkprince/crud-employee-go/configs"
	"github.com/anggadarkprince/crud-employee-go/database"
	"github.com/anggadarkprince/crud-employee-go/middlewares"
	"github.com/anggadarkprince/crud-employee-go/pkg/logger"
	"github.com/anggadarkprince/crud-employee-go/routes"
	"github.com/anggadarkprince/crud-employee-go/utilities"
	"github.com/anggadarkprince/crud-employee-go/pkg/validation"
)

func main() {
	if _, err := configs.Load(); err != nil {
        log.Fatal("Failed to load config:", err)
    }

	logger.Initialize()
	
	utilities.InitTemplates()

	validation.Init()

	db := database.InitDatabase()
	
	server := http.NewServeMux()

	fs := http.FileServer(http.Dir("./uploads"))
    server.Handle("/statics/", http.StripPrefix("/statics/", fs))
	
	server.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/favicon.ico")
	})
	server.HandleFunc("GET /statics/img/no-avatar.png", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/img/no-avatar.png")
	})
	routes.MapRoutes(server, db)

	port := configs.Get().App.Port
	portStr := strconv.Itoa(int(port))

	http.ListenAndServe(":" + portStr, middlewares.MethodOverride(server))
}
