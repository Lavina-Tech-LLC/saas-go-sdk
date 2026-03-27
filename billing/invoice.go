package billing

import "context"

// GetCustomerInvoices returns invoices for a specific customer.
func (s *Service) GetCustomerInvoices(ctx context.Context, customerID string) ([]Invoice, error) {
	var result []Invoice
	err := s.transport.Request(ctx, "GET", "/billing/customers/"+customerID+"/invoices", nil, &result)
	return result, err
}
