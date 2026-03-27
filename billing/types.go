package billing

import "time"

// Customer represents a billing customer.
type Customer struct {
	ID               string    `json:"id"`
	ProjectID        string    `json:"projectId"`
	Email            string    `json:"email"`
	Name             string    `json:"name,omitempty"`
	StripeCustomerID string    `json:"stripeCustomerId,omitempty"`
	BalanceCents     int       `json:"balanceCents"`
	Metadata         string    `json:"metadata,omitempty"`
	TaxExempt        bool      `json:"taxExempt"`
	TaxID            string    `json:"taxId,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

// CreateCustomerParams are the parameters for CreateCustomer.
type CreateCustomerParams struct {
	Email    string `json:"email"`
	Name     string `json:"name,omitempty"`
	Metadata string `json:"metadata,omitempty"`
}

// UpdateCustomerParams are the parameters for UpdateCustomer.
type UpdateCustomerParams struct {
	Email    string `json:"email,omitempty"`
	Name     string `json:"name,omitempty"`
	Metadata string `json:"metadata,omitempty"`
}

// Subscription represents a customer's subscription.
type Subscription struct {
	ID                   string     `json:"id"`
	CustomerID           string     `json:"customerId"`
	PlanID               string     `json:"planId"`
	ProjectID            string     `json:"projectId"`
	Status               string     `json:"status"` // trialing, active, past_due, paused, canceled
	StripeSubscriptionID string     `json:"stripeSubscriptionId,omitempty"`
	CancelAtPeriodEnd    bool       `json:"cancelAtPeriodEnd"`
	TrialEnd             *time.Time `json:"trialEnd,omitempty"`
	CurrentPeriodStart   time.Time  `json:"currentPeriodStart"`
	CurrentPeriodEnd     time.Time  `json:"currentPeriodEnd"`
	CanceledAt           *time.Time `json:"canceledAt,omitempty"`
	CreatedAt            time.Time  `json:"createdAt"`
}

// SubscribeParams are the parameters for Subscribe.
type SubscribeParams struct {
	PlanID string `json:"planId"`
}

// ChangePlanParams are the parameters for ChangePlan.
type ChangePlanParams struct {
	PlanID string `json:"planId"`
}

// CancelResult is returned by CancelSubscription.
type CancelResult struct {
	CanceledAtPeriodEnd bool `json:"canceledAtPeriodEnd"`
}

// UsageEventParams are the parameters for IngestUsageEvent.
type UsageEventParams struct {
	CustomerID     string  `json:"customerId"`
	Metric         string  `json:"metric"`
	Quantity       float64 `json:"quantity"`
	Timestamp      string  `json:"timestamp,omitempty"`
	IdempotencyKey string  `json:"idempotencyKey,omitempty"`
}

// UsageEventResult is returned by IngestUsageEvent.
type UsageEventResult struct {
	ID       string `json:"id"`
	Ingested bool   `json:"ingested"`
}

// UsageSummary represents aggregated usage for a metric.
type UsageSummary struct {
	Metric string  `json:"metric"`
	Total  float64 `json:"total"`
}

// Invoice represents a billing invoice.
type Invoice struct {
	ID              string     `json:"id"`
	ProjectID       string     `json:"projectId"`
	CustomerID      string     `json:"customerId"`
	SubscriptionID  *string    `json:"subscriptionId,omitempty"`
	AmountCents     int        `json:"amountCents"`
	Status          string     `json:"status"` // draft, open, paid, void, uncollectible
	StripeInvoiceID string     `json:"stripeInvoiceId,omitempty"`
	PdfURL          string     `json:"pdfUrl,omitempty"`
	DueDate         *time.Time `json:"dueDate,omitempty"`
	PaidAt          *time.Time `json:"paidAt,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
}

// CreatePortalTokenParams are the parameters for CreatePortalToken.
type CreatePortalTokenParams struct {
	CustomerID string `json:"customerId"`
	ExpiresIn  int    `json:"expiresIn,omitempty"` // seconds, default 3600
}

// PortalTokenResult is returned by CreatePortalToken.
type PortalTokenResult struct {
	PortalToken string    `json:"portalToken"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

// ApplyCouponParams are the parameters for ApplyCoupon.
type ApplyCouponParams struct {
	Code string `json:"code"`
}

// ApplyCouponResult is returned by ApplyCoupon.
type ApplyCouponResult struct {
	Applied      bool   `json:"applied"`
	DiscountType string `json:"discountType"`
	Amount       int    `json:"amount"`
	Duration     string `json:"duration"`
}
