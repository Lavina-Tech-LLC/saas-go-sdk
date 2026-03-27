package report

import "context"

// ExecuteQuery runs a natural language or direct SQL query.
func (s *Service) ExecuteQuery(ctx context.Context, params *QueryParams) (*QueryResult, error) {
	result := &QueryResult{}
	err := s.transport.Request(ctx, "POST", "/reports/query", params, result)
	return result, err
}
