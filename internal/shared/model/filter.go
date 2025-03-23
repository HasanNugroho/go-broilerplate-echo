package shared

type PaginationFilter struct {
	Limit  int    `form:"limit" json:"limit"`
	Page   int    `form:"page" json:"page"`
	Sort   string `form:"sort" json:"sort"`
	Search string `form:"search" json:"search"`
}
