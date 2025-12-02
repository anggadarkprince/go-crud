package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/anggadarkprince/crud-employee-go/dto"
	"github.com/anggadarkprince/crud-employee-go/services"
	"github.com/anggadarkprince/crud-employee-go/utilities"
)

type Employee struct {
	Id int
	Name string
	Email sql.NullString
	TaxNumber sql.NullString
	Gender sql.NullString
	HiredDate sql.NullTime
	Address sql.NullString
	Status sql.NullString
	TotalAllowance int
}

type EmployeeAllowance struct {
	Id int
	EmployeeId int
	Allowance string
}

type EmployeeController struct {
	employeeService *services.EmployeeService
	employeeAllowanceService *services.EmployeeAllowanceService
}

func NewEmployeeController(
	employeeService *services.EmployeeService,
	employeeAllowanceService *services.EmployeeAllowanceService,
) *EmployeeController {
	return &EmployeeController{
		employeeService: employeeService,
		employeeAllowanceService: employeeAllowanceService,
	}
}

func (controller *EmployeeController) Index(w http.ResponseWriter, r *http.Request) error {
	employees, err := controller.employeeService.GetAll(r.Context())
	if err != nil {
        return err
    }

	data := utilities.Compact(
		"employees", employees,
	)

	return utilities.Render(w, r, "employees/index.html", data)
}

func (c *EmployeeController) Create(w http.ResponseWriter, r *http.Request) error {
	return utilities.Render(w, r, "employees/create.html", nil)
}

func (c *EmployeeController) Store(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
        return err
    }

	allowances := r.Form["allowances"]
	data := &dto.CreateEmployeeRequest{
		Name: r.FormValue("name"),
		Email: r.FormValue("email"),
		TaxNumber: r.FormValue("tax_number"),
		Gender: r.FormValue("gender"),
		HiredDate: r.FormValue("hired_date"),
		Address: r.FormValue("address"),
		Status: r.FormValue("status"),
		Allowances: allowances,
	}	
	_, err := c.employeeService.Store(r.Context(), data)
	if err != nil {
        return err
    }
	
	http.Redirect(w, r, "/employees", http.StatusSeeOther)
	return nil
}

func (c *EmployeeController) View(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	employeeId, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
        return err
    }
	employee, err := c.employeeService.GetById(r.Context(), int(employeeId))
	if err != nil {
        return err
    }
	employeeAllowances, err := c.employeeAllowanceService.GetByEmployeeId(r.Context(), int(employeeId))
	if err != nil {
        return err
    }
	
	data := utilities.Compact(
		"employee", employee,
		"employeeAllowances", employeeAllowances,
	)
	return utilities.Render(w, r, "employees/view.html", data)
}

func (c *EmployeeController) Edit(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	employeeId, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
        return err
    }
	employee, err := c.employeeService.GetById(r.Context(), int(employeeId))
	if err != nil {
        return err
    }
	employeeAllowances, err := c.employeeAllowanceService.GetByEmployeeId(r.Context(), int(employeeId))
	if err != nil {
        return err
    }
	
	data := utilities.Compact(
		"employee", employee,
		"employeeAllowances", employeeAllowances,
	)
	return utilities.Render(w, r, "employees/edit.html", data)
}

func (c *EmployeeController) Update(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	employeeId, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
        return err
    }

	if err := r.ParseForm(); err != nil {
        return err
    }

	allowances := r.Form["allowances"]
	data := &dto.UpdateEmployeeRequest{
		Id: int(employeeId),
		Name: r.FormValue("name"),
		Email: r.FormValue("email"),
		TaxNumber: r.FormValue("tax_number"),
		Gender: r.FormValue("gender"),
		HiredDate: r.FormValue("hired_date"),
		Address: r.FormValue("address"),
		Status: r.FormValue("status"),
		Allowances: allowances,
	}	
	_, err = c.employeeService.Update(r.Context(), data)
	if err != nil {
        return err
    }
	
	http.Redirect(w, r, "/employees", http.StatusSeeOther)
	return nil
}

func (c *EmployeeController) Delete(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	employeeId, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
        return err
    }
	err = c.employeeService.Destroy(r.Context(), int(employeeId))
	if err != nil {
        return err
    }
	http.Redirect(w, r, "/employees", http.StatusSeeOther)
	return nil
}