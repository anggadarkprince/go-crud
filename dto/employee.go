package dto

type CreateEmployeeRequest struct {
    Name string `form:"name" validate:"required"`
    Email string `form:"email" validate:"required,email"`
    TaxNumber string `form:"tax_number" validate:"required"`
    Gender string `form:"gender" validate:"required,gender"`
    HiredDate string `form:"hired_date" validate:"required"`
    Address string `form:"address" validate:"required"`
    Status string `form:"status" validate:"required"`
    Allowances []string `form:"allowances" validate:"required"`
}

type UpdateEmployeeRequest struct {
    Id int `validate:"required,number,numeric,gt=0"`
    Name string `form:"name" validate:"required"`
    Email string `form:"email" validate:"required,email"`
    TaxNumber string `form:"tax_number" validate:"required"`
    Gender string `form:"gender" validate:"required,gender"`
    HiredDate string `form:"hired_date" validate:"required"`
    Address string `form:"address" validate:"required"`
    Status string `form:"status" validate:"required"`
    Allowances []string `form:"allowances" validate:"required"`
}
