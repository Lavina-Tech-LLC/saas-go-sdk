package report

import (
	"fmt"
	"net/url"
)

// ListParams holds offset-based pagination parameters for list endpoints.
type ListParams struct {
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"perPage,omitempty"`
	Sort    string `json:"sort,omitempty"`
	Order   string `json:"order,omitempty"`
	Search  string `json:"search,omitempty"`
}

// QueryString converts ListParams to URL query parameters.
func (p *ListParams) QueryString() string {
	if p == nil {
		return ""
	}
	v := url.Values{}
	if p.Page > 0 {
		v.Set("page", fmt.Sprintf("%d", p.Page))
	}
	if p.PerPage > 0 {
		v.Set("perPage", fmt.Sprintf("%d", p.PerPage))
	}
	if p.Sort != "" {
		v.Set("sort", p.Sort)
	}
	if p.Order != "" {
		v.Set("order", p.Order)
	}
	if p.Search != "" {
		v.Set("search", p.Search)
	}
	if encoded := v.Encode(); encoded != "" {
		return "?" + encoded
	}
	return ""
}

// OffsetPage is the standard paginated response envelope.
type OffsetPage[T any] struct {
	Data []T        `json:"data"`
	Meta OffsetMeta `json:"meta"`
}

// OffsetMeta holds pagination metadata.
type OffsetMeta struct {
	Page       int `json:"page"`
	PerPage    int `json:"perPage"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}
