package controllers

import (
	"database/sql"
	"net/http"

	"github.com/anggadarkprince/crud-employee-go/utilities"
)

type DashboardController struct {
	db *sql.DB
}

func NewDashboardController(db *sql.DB) *DashboardController {
	return &DashboardController{db: db}
}

func (c *DashboardController) Index(w http.ResponseWriter, r *http.Request) {
	err := utilities.Render(w, r, "dashboard/index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}