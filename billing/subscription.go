package billing

import "context"

// Subscribe creates a new subscription for a customer.
func (s *Service) Subscribe(ctx context.Context, customerID string, params *SubscribeParams) (*Subscription, error) {
	result := &Subscription{}
	err := s.transport.Request(ctx, "POST", "/billing/customers/"+customerID+"/subscriptions", params, result)
	return result, err
}

// ChangePlan switches a customer's subscription to a different plan.
func (s *Service) ChangePlan(ctx context.Context, customerID string, params *ChangePlanParams) (*Subscription, error) {
	result := &Subscription{}
	err := s.transport.Request(ctx, "PATCH", "/billing/customers/"+customerID+"/subscriptions", params, result)
	return result, err
}

// CancelSubscription marks a subscription to cancel at period end.
func (s *Service) CancelSubscription(ctx context.Context, customerID string) (*CancelResult, error) {
	result := &CancelResult{}
	err := s.transport.Request(ctx, "DELETE", "/billing/customers/"+customerID+"/subscriptions", nil, result)
	return result, err
}
