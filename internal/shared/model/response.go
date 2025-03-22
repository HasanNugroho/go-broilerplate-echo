package shared

// Response format untuk sukses dan error
type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Pagination struct {
	Limit      int   `json:"limit"`
	Page       int   `json:"page"`
	TotalItems int64 `json:"total_items"`
	TotalPages int   `json:"total_pages"`
}

type DataWithPagination struct {
	Items  interface{} `json:"items"`
	Paging Pagination  `json:"paging"`
}
