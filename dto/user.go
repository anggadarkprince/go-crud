package dto

import "mime/multipart"

type UpdateAccountRequest struct {
    Id int `validate:"required,number,numeric,gt=0"`
    Name string `form:"name" validate:"required,min=3,max=50"`
    Username string `form:"username" validate:"required,username,min=3,max=20"`
    Email string `form:"email" validate:"required,email,min=3,max=30"`
	AvatarFile *multipart.FileHeader `form:"avatar"`
	Avatar string `validate:"avatar"`
    CurrentPassword string `form:"current_password" validate:"required,min=3,max=20"`
    Password string `form:"password" validate:"max=20"`
    PasswordConfirmation string `form:"password_confirmation" validate:"eqfield=Password"`
}