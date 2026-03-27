package report

import "context"

// ListQueries returns all saved queries for the project.
func (s *Service) ListQueries(ctx context.Context, params *ListParams) (*OffsetPage[SavedQuery], error) {
	result := &OffsetPage[SavedQuery]{}
	path := "/reports/queries" + params.QueryString()
	err := s.transport.Request(ctx, "GET", path, nil, result)
	return result, err
}

// SaveQuery persists a named query.
func (s *Service) SaveQuery(ctx context.Context, params *SaveQueryParams) (*SavedQuery, error) {
	result := &SavedQuery{}
	err := s.transport.Request(ctx, "POST", "/reports/queries", params, result)
	return result, err
}

// UpdateQuery updates a saved query's name or chart type.
func (s *Service) UpdateQuery(ctx context.Context, queryID string, params *UpdateQueryParams) (*SavedQuery, error) {
	result := &SavedQuery{}
	err := s.transport.Request(ctx, "PATCH", "/reports/queries/"+queryID, params, result)
	return result, err
}

// DeleteQuery removes a saved query.
func (s *Service) DeleteQuery(ctx context.Context, queryID string) (*DeleteResult, error) {
	result := &DeleteResult{}
	err := s.transport.Request(ctx, "DELETE", "/reports/queries/"+queryID, nil, result)
	return result, err
}
