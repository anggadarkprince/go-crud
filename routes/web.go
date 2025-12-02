package routes

import (
	"database/sql"
	"net/http"

	"github.com/anggadarkprince/crud-employee-go/controllers"
)

func MapRoutes(server *http.ServeMux, db *sql.DB) {
	dashboardController := controllers.NewDashboardController(db)
	server.HandleFunc("/", dashboardController.Index)

	employeeController := controllers.NewEmployeeController(db)
	server.HandleFunc("/employees", employeeController.Index)
	server.HandleFunc("/employees/create", employeeController.Create)
	server.HandleFunc("/employees/store", employeeController.Store)
	server.HandleFunc("/employees/{id}", employeeController.View)
	server.HandleFunc("/employees/{id}/edit", employeeController.Edit)
	server.HandleFunc("/employees/{id}/update", employeeController.Update)
	server.HandleFunc("/employees/{id}/delete", employeeController.Delete)
}