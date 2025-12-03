package routes

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/anggadarkprince/crud-employee-go/controllers"
	"github.com/anggadarkprince/crud-employee-go/exceptions"
	"github.com/anggadarkprince/crud-employee-go/middlewares"
	"github.com/anggadarkprince/crud-employee-go/repositories"
	"github.com/anggadarkprince/crud-employee-go/services"
	"github.com/anggadarkprince/crud-employee-go/utilities/session"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error
func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    err := h(w, r)
    if err != nil {
		referer := r.Header.Get("referer")
		accept := r.Header.Get("accept")
		method := r.Method
		fmt.Println("Error", err.Error())
		if (method != "GET" && strings.Contains(accept, "text/html") && referer != "") {
			oldInput := session.ParseFormInput(r)
			message := err.Error()
			var appErr *exceptions.AppError
			if errors.As(err, &appErr) {
				message = appErr.Message
			}
			session.FlashWithInput(w, "danger", message, oldInput)
			http.Redirect(w, r, referer, http.StatusSeeOther)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func MapRoutes(server *http.ServeMux, db *sql.DB) {
	userRepository := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepository)
	authController := controllers.NewAuthController(authService)

	auth := &middlewares.Auth{
		UserRepository: userRepository,
		SecretKey: os.Getenv("JWT_SECRET"),
	}

	server.Handle("/login", auth.GuestMiddleware(HandlerFunc(authController.Index)))
	server.Handle("/authenticate", auth.GuestMiddleware(HandlerFunc(authController.Login)))
	server.Handle("/logout", auth.AuthMiddleware(HandlerFunc(authController.Logout)))

	dashboardRepository := repositories.NewDashboardRepository(db)
	dashboardService := services.NewDashboardService(dashboardRepository)
	dashboardController := controllers.NewDashboardController(dashboardService)
	server.Handle("/", auth.AuthMiddleware(HandlerFunc(dashboardController.Index)))

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
	employeeController := controllers.NewEmployeeController(employeeService, employeeAllowanceService)
	server.Handle("/employees", auth.AuthMiddleware(HandlerFunc(employeeController.Index)))
	server.Handle("/employees/create", auth.AuthMiddleware(HandlerFunc(employeeController.Create)))
	server.Handle("/employees/store", auth.AuthMiddleware(HandlerFunc(employeeController.Store)))
	server.Handle("/employees/{id}", auth.AuthMiddleware(HandlerFunc(employeeController.View)))
	server.Handle("/employees/{id}/edit", auth.AuthMiddleware(HandlerFunc(employeeController.Edit)))
	server.Handle("/employees/{id}/update", auth.AuthMiddleware(HandlerFunc(employeeController.Update)))
	server.Handle("/employees/{id}/delete", auth.AuthMiddleware(HandlerFunc(employeeController.Delete)))
}