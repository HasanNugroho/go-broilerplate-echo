package utils

import (
	"math"

	shared "github.com/HasanNugroho/starter-golang/internal/shared/model"
	"gorm.io/gorm"
)

// Paginate menambahkan pagination ke query GORM
func Paginate(filter *shared.PaginationFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Set default nilai jika kosong
		if filter.Page == 0 {
			filter.Page = 1
		}
		if filter.Limit == 0 {
			filter.Limit = 10
		}
		if filter.Sort == "" {
			filter.Sort = "id desc"
		}

		return db.Offset((filter.Page - 1) * filter.Limit).Limit(filter.Limit).Order(filter.Sort)
	}
}

func BuildPagination(filter *shared.PaginationFilter, totalItems int64) shared.Pagination {
	// Hitung total halaman
	totalPages := int(math.Ceil(float64(totalItems) / float64(filter.Limit)))

	// Buat response dengan pagination
	return shared.Pagination{
		Limit:      filter.Limit,
		Page:       filter.Page,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
