package controllers

import (
	"errors"
	"net/http"
	"os"

	appErrors "github.com/anggadarkprince/crud-employee-go/errors"
	"github.com/anggadarkprince/crud-employee-go/services"
	"github.com/anggadarkprince/crud-employee-go/utilities"
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

	user, err := controller.authService.Authenticate(r.Context(), username, password)
	if err != nil {
		switch {
        case errors.Is(err, appErrors.ErrUserNotFound):
            http.Error(w, "User not found", http.StatusNotFound)
			return nil
        case errors.Is(err, appErrors.ErrWrongPassword):
            http.Error(w, "Username or password wrong", http.StatusUnauthorized)
			return nil
        case errors.Is(err, appErrors.ErrUserInactive):
            http.Error(w, "User is PENDING or SUSPENDED", http.StatusUnauthorized)
			return nil
        default:
            http.Error(w, "Can't logged you in", http.StatusInternalServerError)
			return nil
        }
	}
	
	if user != nil {
		authToken, err := controller.authService.GenerateAuthToken(user.Id)
		if err != nil {
			return err
		}

		var cookieName = os.Getenv("COOKIE_NAME")
		if cookieName == "" {
			cookieName = "auth_token"
		}

		cookie := http.Cookie{
			Name: cookieName,
			Value: authToken,
			Path: "/",
			HttpOnly: true, // cannot be accessed by JS (secure)
			Secure: false, // set to true in HTTPS
			SameSite: http.SameSiteLaxMode,
			MaxAge: 86400, // 1 day
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

    http.Redirect(w, r, "/login", http.StatusSeeOther)
	return nil
}