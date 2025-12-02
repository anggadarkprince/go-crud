package repositories

import (
	"context"
	"database/sql"
)

type DashboardStatistics struct {
    Total    int
    Active   int
    Inactive int
    Pending  int
}

type DashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (repository *DashboardRepository) GetStatistics(ctx context.Context) (*DashboardStatistics, error) {
	var stats DashboardStatistics

	query := `
		SELECT
			COUNT(*) AS total_employees,
			COUNT(IF(status = 'ACTIVE', 1, NULL)) AS active_employees,
			COUNT(IF(status = 'INACTIVE', 1, NULL)) AS inactive_employees,
			COUNT(IF(status = 'PENDING', 1, NULL)) AS pending_employees
		FROM employees
	`
	err := repository.db.
		QueryRowContext(ctx, query).
		Scan(&stats.Total, &stats.Active, &stats.Inactive, &stats.Pending)
	
	if err != nil {
        return nil, err
    }

	return &stats, nil
}