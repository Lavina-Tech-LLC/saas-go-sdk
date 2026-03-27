package billing

import "context"

// IngestUsageEvent records a metered usage event.
func (s *Service) IngestUsageEvent(ctx context.Context, params *UsageEventParams) (*UsageEventResult, error) {
	result := &UsageEventResult{}
	err := s.transport.Request(ctx, "POST", "/billing/usage", params, result)
	return result, err
}

// GetCurrentUsage returns aggregated usage for the current billing period.
func (s *Service) GetCurrentUsage(ctx context.Context, customerID string) ([]UsageSummary, error) {
	var result []UsageSummary
	err := s.transport.Request(ctx, "GET", "/billing/customers/"+customerID+"/usage", nil, &result)
	return result, err
}
