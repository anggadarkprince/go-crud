package routes

import (
	"database/sql"
	"net/http"

	"github.com/anggadarkprince/crud-employee-go/controllers"
	"github.com/anggadarkprince/crud-employee-go/repositories"
	"github.com/anggadarkprince/crud-employee-go/services"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error
func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    err := h(w, r)
    if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func MapRoutes(server *http.ServeMux, db *sql.DB) {
	dashboardRepository := repositories.NewDashboardRepository(db)
	dashboardService := services.NewDashboardService(dashboardRepository)
	dashboardController := controllers.NewDashboardController(dashboardService)
	server.Handle("/", HandlerFunc(dashboardController.Index))

	employeeRepository := repositories.NewEmployeeRepository(db)
	employeeAllowanceRepository := repositories.NewEmployeeAllowanceRepository(db)
	employeeService := services.NewEmployeeService(
		employeeRepository,
		employeeAllowanceRepository,
		db,
	)
	employeeAllowanceService := services.NewEmployeeAllowanceService(
		employeeAllowanceRepository,
	)
	employeeController := controllers.NewEmployeeController(db, employeeService, employeeAllowanceService)
	server.Handle("/employees", HandlerFunc(employeeController.Index))
	server.Handle("/employees/create", HandlerFunc(employeeController.Create))
	server.Handle("/employees/store", HandlerFunc(employeeController.Store))
	server.Handle("/employees/{id}", HandlerFunc(employeeController.View))
	server.Handle("/employees/{id}/edit", HandlerFunc(employeeController.Edit))
	server.Handle("/employees/{id}/update", HandlerFunc(employeeController.Update))
	server.Handle("/employees/{id}/delete", HandlerFunc(employeeController.Delete))
}