package controllers

import (
	"net/http"

	"github.com/anggadarkprince/crud-employee-go/utilities"
)

func NewDashboardController() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := utilities.Render(w, r, "dashboard/index.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}