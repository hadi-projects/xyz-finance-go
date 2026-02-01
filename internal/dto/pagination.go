package dto

// PaginationRequest for incoming pagination params
type PaginationRequest struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

// PaginationMeta for response pagination metadata
type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// PaginatedResponse wraps data with pagination metadata
type PaginatedResponse struct {
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// GetOffset calculates SQL offset from page and limit
func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

// SetDefaults sets default values if not provided
func (p *PaginationRequest) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 20
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
}

// NewPaginationMeta creates pagination metadata from total count
func NewPaginationMeta(page, limit int, total int64) PaginationMeta {
	totalPages := (total + int64(limit) - 1) / int64(limit) // ceiling division
	return PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    int64(page) < totalPages,
		HasPrev:    page > 1,
	}
}

// NewPaginatedResponse creates a paginated response
func NewPaginatedResponse(data interface{}, page, limit int, total int64) PaginatedResponse {
	return PaginatedResponse{
		Data:       data,
		Pagination: NewPaginationMeta(page, limit, total),
	}
}
