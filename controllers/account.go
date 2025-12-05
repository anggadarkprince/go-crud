package controllers

import (
	"net/http"

	"github.com/anggadarkprince/crud-employee-go/dto"
	"github.com/anggadarkprince/crud-employee-go/middlewares"
	"github.com/anggadarkprince/crud-employee-go/pkg/session"
	"github.com/anggadarkprince/crud-employee-go/pkg/validation"
	"github.com/anggadarkprince/crud-employee-go/services"
	"github.com/anggadarkprince/crud-employee-go/utilities"
)

type AccountController struct {
	userService *services.UserService
}

func NewAccountController(userService *services.UserService) *AccountController {
	return &AccountController{userService: userService}
}

func (controller *AccountController) Index(w http.ResponseWriter, r *http.Request) error {
	user := middlewares.GetUser(r)
	data := utilities.Compact("user", user)
	return utilities.Render(w, r, "account/index.html", data)
}

func (controller *AccountController) Update(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return err
	}
	user := middlewares.GetUser(r)
	
	data := &dto.UpdateAccountRequest{
		Id:         int(user.Id),
		Name:       r.FormValue("name"),
		Username: r.FormValue("username"),
		Email:      r.FormValue("email"),
		CurrentPassword:  r.FormValue("current_password"),
		Password:     r.FormValue("password"),
		PasswordConfirmation:  r.FormValue("password_confirmation"),
	}
	file, header, err := r.FormFile("avatar")
    if err == nil {
        data.AvatarFile = header
        file.Close()
    }

	err = validation.Validator.Struct(data)
	if err != nil {
		return err
	}

	_, err = controller.userService.UpdateAccount(r.Context(), data)
	if err != nil {
		return err
	}
	
	session.Flash(w, "success", "Account successfully updated")

	http.Redirect(w, r, "/account", http.StatusSeeOther)
	return nil
}