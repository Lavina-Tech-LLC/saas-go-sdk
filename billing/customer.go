package billing

import "context"

// CreateCustomer creates a billing customer, optionally syncing to Stripe.
func (s *Service) CreateCustomer(ctx context.Context, params *CreateCustomerParams) (*Customer, error) {
	result := &Customer{}
	err := s.transport.Request(ctx, "POST", "/billing/customers", params, result)
	return result, err
}

// GetCustomer returns a billing customer by ID.
func (s *Service) GetCustomer(ctx context.Context, customerID string) (*Customer, error) {
	result := &Customer{}
	err := s.transport.Request(ctx, "GET", "/billing/customers/"+customerID, nil, result)
	return result, err
}

// UpdateCustomer updates mutable fields on a billing customer.
func (s *Service) UpdateCustomer(ctx context.Context, customerID string, params *UpdateCustomerParams) (*Customer, error) {
	result := &Customer{}
	err := s.transport.Request(ctx, "PATCH", "/billing/customers/"+customerID, params, result)
	return result, err
}
