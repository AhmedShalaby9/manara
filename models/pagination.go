package models

type Pagination struct {
	TotalPages int `json:"total_pages"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalItems int `json:"total_items"`
}
type PaginationParams struct {
	Page    int
	PerPage int
}
