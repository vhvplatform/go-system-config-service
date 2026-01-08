package domain

// Common request/response structures

// PaginationRequest represents pagination parameters
type PaginationRequest struct {
	Page    int `form:"page" json:"page"`
	PerPage int `form:"per_page" json:"per_page"`
}

// SetDefaults sets default values for pagination
func (p *PaginationRequest) SetDefaults() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage < 1 {
		p.PerPage = 30
	}
	if p.PerPage > 100 {
		p.PerPage = 100
	}
}

// PaginationResponse represents pagination metadata
type PaginationResponse struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	TotalPages int   `json:"total_pages"`
	TotalItems int64 `json:"total_items"`
}

// ListResponse represents a generic list response
type ListResponse struct {
	Data       interface{}        `json:"data"`
	Pagination PaginationResponse `json:"pagination,omitempty"`
}

// CreateResponse represents a generic create response
type CreateResponse struct {
	ID string `json:"id"`
}

// UpdateResponse represents a generic update response
type UpdateResponse struct {
	ID string `json:"id"`
}

// DeleteResponse represents a generic delete response
type DeleteResponse struct {
	ID string `json:"id"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
	Version string `json:"version"`
}
