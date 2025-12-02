package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

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
	db *sql.DB
}

func NewEmployeeController(db *sql.DB) *EmployeeController {
	return &EmployeeController{db: db}
}

func (c *EmployeeController) Index(w http.ResponseWriter, r *http.Request) {
	rows, err := c.db.Query(`
		SELECT 
			id, name, email, tax_number, gender, hired_date, address, status, 
			(SELECT COUNT(*) FROM employee_allowances WHERE employee_id = employees.id) AS total_allowance 
		FROM employees
		ORDER BY id DESC
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var employee Employee

		err = rows.Scan(
			&employee.Id,
			&employee.Name,
			&employee.Email,
			&employee.TaxNumber,
			&employee.Gender,
			&employee.HiredDate,
			&employee.Address,
			&employee.Status,
			&employee.TotalAllowance,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		employees = append(employees, employee)
	}

	data := make(map[string]any)
	data["employees"] = employees

	utilities.Render(w, r, "employees/index.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *EmployeeController) Create(w http.ResponseWriter, r *http.Request) {
	err := utilities.Render(w, r, "employees/create.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *EmployeeController) Store(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := r.FormValue("name")
	email := r.FormValue("email")
	taxNumber := r.FormValue("tax_number")
	gender := r.FormValue("gender")
	hiredDate := r.FormValue("hired_date")
	address := r.FormValue("address")
	status := r.FormValue("status")
	allowances := r.Form["allowances"]

	var hiredDateValue *string
	if hiredDate != "" {
		hiredDateValue = &hiredDate
	} else {
		hiredDateValue = nil
	}

	fmt.Println(name, email, taxNumber, gender, hiredDateValue, address, status, allowances)

	tx, err := c.db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	query := "INSERT INTO employees(name, email, tax_number, gender, hired_date, address, status) VALUES(?, ?, ?, ?, ?, ?, ?)"
	result, err := tx.Exec(query, name, email, taxNumber, gender, hiredDateValue, address, status)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	employeeId, _ := result.LastInsertId()

	statement, _ := tx.Prepare("INSERT INTO employee_allowances(employee_id, allowance) VALUES(?, ?)")
	defer statement.Close()

	for _, item := range allowances {
		_, err := statement.Exec(employeeId, item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	http.Redirect(w, r, "/employees", http.StatusMovedPermanently)
}

func (c *EmployeeController) View(w http.ResponseWriter, r *http.Request) {
	employeeId := r.PathValue("id")
	row := c.db.QueryRow("SELECT id, name, email, tax_number, gender, hired_date, address, status FROM employees WHERE id = ?", employeeId)
	if row.Err() != nil {
		http.Error(w, row.Err().Error(), http.StatusInternalServerError)
		return
	}
	var employee Employee
	err := row.Scan(
		&employee.Id,
		&employee.Name,
		&employee.Email,
		&employee.TaxNumber,
		&employee.Gender,
		&employee.HiredDate,
		&employee.Address,
		&employee.Status,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(employeeId)		
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	employee.Id = id

	rows, err := c.db.Query("SELECT id, employee_id, allowance FROM employee_allowances WHERE employee_id = ?", employeeId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var employeeAllowances []EmployeeAllowance
	for rows.Next() {
		var employeeAllowance EmployeeAllowance

		err = rows.Scan(
			&employeeAllowance.Id,
			&employeeAllowance.EmployeeId,
			&employeeAllowance.Allowance,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		employeeAllowances = append(employeeAllowances, employeeAllowance)
	}

	data := make(map[string]any)
	data["employee"] = employee
	data["employeeAllowances"] = employeeAllowances
	
	err = utilities.Render(w, r, "employees/view.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *EmployeeController) Edit(w http.ResponseWriter, r *http.Request) {
	employeeId := r.PathValue("id")
	row := c.db.QueryRow("SELECT id, name, email, tax_number, gender, hired_date, address, status FROM employees WHERE id = ?", employeeId)
	if row.Err() != nil {
		http.Error(w, row.Err().Error(), http.StatusInternalServerError)
		return
	}
	var employee Employee
	err := row.Scan(
		&employee.Id,
		&employee.Name,
		&employee.Email,
		&employee.TaxNumber,
		&employee.Gender,
		&employee.HiredDate,
		&employee.Address,
		&employee.Status,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := strconv.Atoi(employeeId)		
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	employee.Id = id

	rows, err := c.db.Query("SELECT id, employee_id, allowance FROM employee_allowances WHERE employee_id = ?", employeeId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var employeeAllowances []EmployeeAllowance
	for rows.Next() {
		var employeeAllowance EmployeeAllowance

		err = rows.Scan(
			&employeeAllowance.Id,
			&employeeAllowance.EmployeeId,
			&employeeAllowance.Allowance,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		employeeAllowances = append(employeeAllowances, employeeAllowance)
	}

	data := make(map[string]any)
	data["employee"] = employee
	data["employeeAllowances"] = employeeAllowances
	
	err = utilities.Render(w, r, "employees/edit.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *EmployeeController) Update(w http.ResponseWriter, r *http.Request) {
	employeeId := r.PathValue("id")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	taxNumber := r.FormValue("tax_number")
	gender := r.FormValue("gender")
	hiredDate := r.FormValue("hired_date")
	address := r.FormValue("address")
	status := r.FormValue("status")
	allowances := r.Form["allowances"]

	var hiredDateValue *string
	if hiredDate != "" {
		hiredDateValue = &hiredDate
	} else {
		hiredDateValue = nil
	}

	fmt.Println(name, email, taxNumber, gender, hiredDateValue, address, status, allowances)

	tx, err := c.db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	query := "UPDATE employees SET name = ?, email = ?, tax_number = ?, gender = ?, hired_date = ?, address = ?, status = ? WHERE id = ?"
	_, err = tx.Exec(query, name, email, taxNumber, gender, hiredDateValue, address, status, employeeId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("DELETE FROM employee_allowances WHERE employee_id = ?", employeeId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	statement, _ := tx.Prepare("INSERT INTO employee_allowances(employee_id, allowance) VALUES(?, ?)")
	defer statement.Close()

	for _, item := range allowances {
		_, err := statement.Exec(employeeId, item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	http.Redirect(w, r, "/employees", http.StatusMovedPermanently)
}

func (c *EmployeeController) Delete(w http.ResponseWriter, r *http.Request) {
	employeeId := r.PathValue("id")
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := c.db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM employees WHERE id = ?", employeeId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("DELETE FROM employee_allowances WHERE employee_id = ?", employeeId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	http.Redirect(w, r, "/employees", http.StatusMovedPermanently)
}