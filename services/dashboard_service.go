package services

import (
	"context"

	"github.com/anggadarkprince/crud-employee-go/repositories"
)

type DashboardService struct {
	dashboardRepository *repositories.DashboardRepository
}

func NewDashboardService(dashboardRepository *repositories.DashboardRepository) *DashboardService {
	return &DashboardService{dashboardRepository: dashboardRepository}
}

func (service *DashboardService) GetStatistics(ctx context.Context) (*repositories.DashboardStatistics, error) {
	return service.dashboardRepository.GetStatistics(ctx)
}
