package controllers

import (
	"net/http"

	"github.com/anggadarkprince/crud-employee-go/services"
	"github.com/anggadarkprince/crud-employee-go/utilities"
)

type DashboardController struct {
	dashboardService *services.DashboardService
}

func NewDashboardController(dashboardService *services.DashboardService) *DashboardController {
	return &DashboardController{dashboardService: dashboardService}
}

func (controller *DashboardController) Index(w http.ResponseWriter, r *http.Request) error {
	stats, err := controller.dashboardService.GetStatistics(r.Context())
	if err != nil {
        return err
    }

	data := utilities.Compact(
		"statistic", stats,
	)

	return utilities.Render(w, r, "dashboard/index.html", data)
}