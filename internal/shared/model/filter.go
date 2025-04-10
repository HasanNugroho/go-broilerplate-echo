package shared

type PaginationFilter struct {
	Limit  int    `form:"limit" json:"limit" query:"limit"`
	Page   int    `form:"page" json:"page" query:"page"`
	Sort   string `form:"sort" json:"sort" query:"sort"`
	Search string `form:"search" json:"search" query:"search"`
}
