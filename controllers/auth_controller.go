package controllers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/anggadarkprince/crud-employee-go/exceptions"
	"github.com/anggadarkprince/crud-employee-go/services"
	"github.com/anggadarkprince/crud-employee-go/utilities"
	"github.com/anggadarkprince/crud-employee-go/utilities/session"
	"github.com/anggadarkprince/crud-employee-go/utilities/validation"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (controller *AuthController) Index(w http.ResponseWriter, r *http.Request) error {
	return utilities.Render(w, r, "auth/login.html", nil)
}

func (controller *AuthController) Login(w http.ResponseWriter, r *http.Request) error {
	username := r.FormValue("username")
	password := r.FormValue("password")
	remember := r.FormValue("remember") == "1"

	validationErrs := make(map[string]string)
	err := validation.Validator.Var(username, "required,min=3,max=50,username")
	if err != nil {
		validationErrs["username"] = "This username should required and valid"
	}
	err = validation.Validator.Var(password, "required")
	if err != nil {
		validationErrs["password"] = "Password is required"
	}
	if len(validationErrs) > 0 {
		return &exceptions.ValidationError{
			Message: "Data is invalid",
			Errors: validationErrs,
		}
	}

	user, err := controller.authService.Authenticate(r.Context(), username, password)
	if err != nil {
		switch {
        case errors.Is(err, exceptions.ErrUserNotFound):
			return &exceptions.AppError{
				Code: 404,
				Message: "User not found",
			}
        case errors.Is(err, exceptions.ErrWrongPassword):
			return &exceptions.AppError{
				Code: 401,
				Message: "Username or password wrong",
			}
        case errors.Is(err, exceptions.ErrUserInactive):
			return &exceptions.AppError{
				Code: 403,
				Message: "User is PENDING or SUSPENDED",
			}
        default:
			return err
        }
	}
	
	if user != nil {
		var hours = 2;
		if remember {
			hours = 24 * 30
		}
		var exp = time.Now().Add(time.Duration(hours) * time.Hour).Unix()
		authToken, err := controller.authService.GenerateAuthToken(user.Id, exp)
		if err != nil {
			return err
		}

		var cookieName = os.Getenv("COOKIE_NAME")
		if cookieName == "" {
			cookieName = "auth_token"
		}
		var maxAge = 7200;
		if remember {
			maxAge = 3600 * 24 * 30
		}

		cookie := http.Cookie{
			Name: cookieName,
			Value: authToken,
			Path: "/",
			HttpOnly: true, // cannot be accessed by JS (secure)
			Secure: false, // set to true in HTTPS
			SameSite: http.SameSiteLaxMode,
			MaxAge: maxAge,
		}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil
	}
	return err
}

func (controller *AuthController) Logout(w http.ResponseWriter, r *http.Request) error {
    var cookieName = os.Getenv("COOKIE_NAME")
	if cookieName == "" {
		cookieName = "auth_token"
	}

    cookie := http.Cookie{
        Name: cookieName,
        Value: "",
        Path: "/",
        MaxAge: -1,
        HttpOnly: true,
        Secure: false, // set true in production (HTTPS)
        SameSite: http.SameSiteLaxMode,
    }
    http.SetCookie(w, &cookie)

	session.Flash(w, "warning", "You are logged out")

    http.Redirect(w, r, "/login", http.StatusSeeOther)
	return nil
}