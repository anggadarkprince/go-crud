package dto

type CreateEmployeeRequest struct {
    Name string `form:"name"`
    Email string `form:"email"`
    TaxNumber string `form:"tax_number"`
    Gender string `form:"gender"`
    HiredDate string `form:"hired_date"`
    Address string `form:"address"`
    Status string `form:"status"`
    Allowances []string `form:"allowances"`
}

type UpdateEmployeeRequest struct {
    Id int
    Name string `form:"name"`
    Email string `form:"email"`
    TaxNumber string `form:"tax_number"`
    Gender string `form:"gender"`
    HiredDate string `form:"hired_date"`
    Address string `form:"address"`
    Status string `form:"status"`
    Allowances []string `form:"allowances"`
}
