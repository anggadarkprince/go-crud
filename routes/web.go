package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/anggadarkprince/crud-employee-go/configs"
	"github.com/anggadarkprince/crud-employee-go/controllers"
	"github.com/anggadarkprince/crud-employee-go/exceptions"
	"github.com/anggadarkprince/crud-employee-go/middlewares"
	"github.com/anggadarkprince/crud-employee-go/pkg/logger"
	"github.com/anggadarkprince/crud-employee-go/pkg/session"
	"github.com/anggadarkprince/crud-employee-go/pkg/validation"
	"github.com/anggadarkprince/crud-employee-go/repositories"
	"github.com/anggadarkprince/crud-employee-go/services"
	"github.com/go-playground/validator/v10"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error
func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    err := h(w, r)
    if err != nil {
		var errorMessage string
		var errorData map[string]string

		// Validation error
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			validationErrors := validation.FormatValidationErrors(validationErrors)
			errorMessage = "Please check the data you provided."
			errorData = validationErrors
		} else if validationErrors, ok := err.(*exceptions.ValidationError); ok {
			errorMessage = validationErrors.Message
			errorData = validationErrors.Errors
		} else {
			logger.LogError("Uncaught exception", err, r)
		}
		
		referer := r.Header.Get("referer")
		accept := r.Header.Get("accept")
		method := r.Method
		if (method != "GET" && strings.Contains(accept, "text/html") && referer != "") {
			oldInput := session.ParseFormInput(r)

			var appErr *exceptions.AppError
			if errors.As(err, &appErr) {
				errorMessage = appErr.Message
			}
			if errorMessage == "" {
				if configs.Get().App.Environment == "production" {
					errorMessage = "Something went wrong"
				} else {
					errorMessage = err.Error()
				}
			}
			
			flashData := session.FlashData{
				"alert": map[string]string{
					"type": "danger",
					"message": errorMessage,
				},
				"old": oldInput,
				"error": errorMessage,
				"errors": errorData,
			}
			session.SetFlash(w, flashData)
			http.Redirect(w, r, referer, http.StatusSeeOther)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func authGroup(auth *middlewares.Auth, routes map[string]http.Handler) map[string]http.Handler {
    grouped := make(map[string]http.Handler)
    for pattern, handler := range routes {
        grouped[pattern] = auth.AuthMiddleware(handler)
    }
    return grouped
}

func guestGroup(auth *middlewares.Auth, routes map[string]http.Handler) map[string]http.Handler {
    grouped := make(map[string]http.Handler)
    for pattern, handler := range routes {
        grouped[pattern] = auth.GuestMiddleware(handler)
    }
    return grouped
}

func registerRoutes(mux *http.ServeMux, routes map[string]http.Handler) {
    for pattern, handler := range routes {
        mux.Handle(pattern, handler)
    }
}

func MapRoutes(server *http.ServeMux, db *sql.DB) {
	userRepository := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepository)
	authController := controllers.NewAuthController(authService)

	auth := &middlewares.Auth{
		UserRepository: userRepository,
		SecretKey: configs.Get().Auth.JwtSecret,
	}

	// Guest routes
	registerRoutes(server, guestGroup(auth, map[string]http.Handler{
		"GET /login": HandlerFunc(authController.Login),
        "POST /login": HandlerFunc(authController.Authenticate),
        "GET /register": HandlerFunc(authController.Register),
        "POST /register": HandlerFunc(authController.RegisterUser),
	}))
	server.Handle("GET /logout", auth.AuthMiddleware(HandlerFunc(authController.Logout)))

	dashboardRepository := repositories.NewDashboardRepository(db)
	dashboardService := services.NewDashboardService(dashboardRepository)
	dashboardController := controllers.NewDashboardController(dashboardService)
	server.Handle("GET /{$}", auth.AuthMiddleware(HandlerFunc(dashboardController.Index)))
	server.Handle("GET /dashboard", auth.AuthMiddleware(HandlerFunc(dashboardController.Index)))

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

	// Auth-protected routes
    registerRoutes(server, authGroup(auth, map[string]http.Handler{
        "GET /employees": HandlerFunc(employeeController.Index),
        "GET /employees/create": HandlerFunc(employeeController.Create),
        "POST /employees": HandlerFunc(employeeController.Store),
        "GET /employees/{id}": HandlerFunc(employeeController.View),
        "GET /employees/{id}/edit": HandlerFunc(employeeController.Edit),
        "PUT /employees/{id}": HandlerFunc(employeeController.Update),
        "DELETE /employees/{id}": HandlerFunc(employeeController.Delete),
    }))
}