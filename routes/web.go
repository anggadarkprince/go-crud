package routes

import (
	"database/sql"
	"net/http"

	"github.com/anggadarkprince/crud-employee-go/controllers"
)

func MapRoutes(server *http.ServeMux, db *sql.DB) {
	server.HandleFunc("/", controllers.NewDashboardController())
	server.HandleFunc("/employees", controllers.NewIndexEmployeeController(db))
	server.HandleFunc("/employees/create", controllers.NewCreateEmployeeController())
	server.HandleFunc("/employees/store", controllers.NewStoreEmployeeController(db))
	server.HandleFunc("/employees/{id}", controllers.NewViewEmployeeController(db))
	server.HandleFunc("/employees/{id}/edit", controllers.NewEditEmployeeController(db))
	server.HandleFunc("/employees/{id}/update", controllers.NewUpdateEmployeeController(db))
	server.HandleFunc("/employees/{id}/delete", controllers.NewDeleteEmployeeController(db))
}