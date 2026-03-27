package report

import "time"

// QueryParams defines the input for executing a query.
type QueryParams struct {
	NaturalLanguage string       `json:"naturalLanguage,omitempty"`
	SQL             string       `json:"sql,omitempty"`
	FilterRules     []FilterRule `json:"filterRules,omitempty"`
}

// FilterRule constrains query results.
type FilterRule struct {
	Table  string `json:"table,omitempty"`
	Column string `json:"column"`
	Op     string `json:"op"` // eq, neq, gt, gte, lt, lte, in, between
	Value  string `json:"value"`
}

// QueryResult is returned by ExecuteQuery.
type QueryResult struct {
	SQL         string                   `json:"sql"`
	Columns     []string                 `json:"columns"`
	Rows        []map[string]interface{} `json:"rows"`
	RowCount    int                      `json:"rowCount"`
	ExecutionMs int                      `json:"executionMs"`
	ChartType   string                   `json:"chartType"`
	Cached      bool                     `json:"cached"`
}

// SavedQuery represents a saved query.
type SavedQuery struct {
	ID              string    `json:"id"`
	ProjectID       string    `json:"projectId"`
	Name            string    `json:"name"`
	NaturalLanguage string    `json:"naturalLanguage,omitempty"`
	GeneratedSQL    string    `json:"generatedSql,omitempty"`
	ChartType       string    `json:"chartType"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// SaveQueryParams are the parameters for SaveQuery.
type SaveQueryParams struct {
	Name            string `json:"name"`
	NaturalLanguage string `json:"naturalLanguage,omitempty"`
	GeneratedSQL    string `json:"generatedSql,omitempty"`
	ChartType       string `json:"chartType,omitempty"`
}

// UpdateQueryParams are the parameters for UpdateQuery.
type UpdateQueryParams struct {
	Name      string `json:"name,omitempty"`
	ChartType string `json:"chartType,omitempty"`
}

// Dashboard represents a report dashboard.
type Dashboard struct {
	ID                     string    `json:"id"`
	ProjectID              string    `json:"projectId"`
	Name                   string    `json:"name"`
	LayoutJSON             string    `json:"layoutJson"`
	IsPublic               bool      `json:"isPublic"`
	RefreshIntervalSeconds int       `json:"refreshIntervalSeconds"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
}

// CreateDashboardParams are the parameters for CreateDashboard.
type CreateDashboardParams struct {
	Name                   string `json:"name"`
	LayoutJSON             string `json:"layoutJson,omitempty"`
	IsPublic               bool   `json:"isPublic,omitempty"`
	RefreshIntervalSeconds int    `json:"refreshIntervalSeconds,omitempty"`
}

// UpdateDashboardParams are the parameters for UpdateDashboard.
type UpdateDashboardParams struct {
	Name                   string `json:"name,omitempty"`
	LayoutJSON             string `json:"layoutJson,omitempty"`
	IsPublic               bool   `json:"isPublic"`
	RefreshIntervalSeconds int    `json:"refreshIntervalSeconds"`
}

// CreateEmbedTokenParams are the parameters for CreateEmbedToken.
type CreateEmbedTokenParams struct {
	DashboardID string       `json:"dashboardId,omitempty"`
	CustomerID  string       `json:"customerId,omitempty"`
	ExpiresIn   int          `json:"expiresIn,omitempty"` // seconds, default 3600
	FilterRules []FilterRule `json:"filterRules,omitempty"`
}

// EmbedTokenResult is returned by CreateEmbedToken.
type EmbedTokenResult struct {
	EmbedToken string    `json:"embedToken"`
	ExpiresAt  time.Time `json:"expiresAt"`
}

// EmbedToken represents an embed token record.
type EmbedToken struct {
	ID          string     `json:"id"`
	ProjectID   string     `json:"projectId"`
	DashboardID *string    `json:"dashboardId,omitempty"`
	CustomerID  string     `json:"customerId,omitempty"`
	FilterRules string     `json:"filterRules"`
	ExpiresAt   time.Time  `json:"expiresAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	RevokedAt   *time.Time `json:"revokedAt,omitempty"`
}

// DeleteResult is returned by delete operations.
type DeleteResult struct {
	Deleted bool `json:"deleted"`
}

// RevokeResult is returned by revoke operations.
type RevokeResult struct {
	Revoked bool `json:"revoked"`
}
