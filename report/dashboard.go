package report

import "context"

// ListDashboards returns all dashboards for the project.
func (s *Service) ListDashboards(ctx context.Context, params *ListParams) (*OffsetPage[Dashboard], error) {
	result := &OffsetPage[Dashboard]{}
	path := "/reports/dashboards" + params.QueryString()
	err := s.transport.Request(ctx, "GET", path, nil, result)
	return result, err
}

// CreateDashboard creates a new dashboard.
func (s *Service) CreateDashboard(ctx context.Context, params *CreateDashboardParams) (*Dashboard, error) {
	result := &Dashboard{}
	err := s.transport.Request(ctx, "POST", "/reports/dashboards", params, result)
	return result, err
}

// GetDashboard returns a single dashboard by ID.
func (s *Service) GetDashboard(ctx context.Context, dashboardID string) (*Dashboard, error) {
	result := &Dashboard{}
	err := s.transport.Request(ctx, "GET", "/reports/dashboards/"+dashboardID, nil, result)
	return result, err
}

// UpdateDashboard updates a dashboard.
func (s *Service) UpdateDashboard(ctx context.Context, dashboardID string, params *UpdateDashboardParams) (*Dashboard, error) {
	result := &Dashboard{}
	err := s.transport.Request(ctx, "PATCH", "/reports/dashboards/"+dashboardID, params, result)
	return result, err
}

// DeleteDashboard removes a dashboard.
func (s *Service) DeleteDashboard(ctx context.Context, dashboardID string) (*DeleteResult, error) {
	result := &DeleteResult{}
	err := s.transport.Request(ctx, "DELETE", "/reports/dashboards/"+dashboardID, nil, result)
	return result, err
}
