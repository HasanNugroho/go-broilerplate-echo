package utils

import (
	"math"

	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
)

func BuildPagination(filter *shared.PaginationFilter, totalItems int64) shared.Pagination {
	// Hitung total halaman
	totalPages := 1
	if totalItems > 1 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(filter.Limit)))
	}

	// Buat response dengan pagination
	return shared.Pagination{
		Limit:      filter.Limit,
		Page:       filter.Page,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
