package dto

type RegisterUserRequest struct {
    Name string `form:"name" validate:"required,min=3,max=50"`
    Username string `form:"username" validate:"required,username,min=3,max=20"`
    Email string `form:"email" validate:"required,email,min=3,max=30"`
    Password string `form:"password" validate:"required,min=3,max=20"`
    PasswordConfirmation string `form:"password_confirmation" validate:"required,eqfield=Password"`
    Agreement string `form:"agreement" validate:"required,oneof=0 1 yes no"`
}