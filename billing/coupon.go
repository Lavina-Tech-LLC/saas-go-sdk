package billing

import "context"

// ApplyCoupon applies a coupon code to a customer's subscription.
func (s *Service) ApplyCoupon(ctx context.Context, customerID string, params *ApplyCouponParams) (*ApplyCouponResult, error) {
	result := &ApplyCouponResult{}
	err := s.transport.Request(ctx, "POST", "/billing/customers/"+customerID+"/coupons", params, result)
	return result, err
}
