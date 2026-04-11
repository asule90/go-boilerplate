package utils

import (
	"math"

	"github.com/gofiber/fiber/v2"
)

// PaginationMeta holds pagination metadata for API responses.
type PaginationMeta struct {
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
	TotalPages   int   `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
}

// CalculatePagination computes pagination metadata.
func CalculatePagination(total int64, page, pageSize int) PaginationMeta {
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages < 1 {
		totalPages = 1
	}
	return PaginationMeta{
		CurrentPage:  page,
		PageSize:     pageSize,
		TotalPages:   totalPages,
		TotalRecords: total,
	}
}

// ParsePaginationParams extracts and validates pagination query parameters.
func ParsePaginationParams(c *fiber.Ctx) (page, pageSize int) {
	page = c.QueryInt("page", 1)
	pageSize = c.QueryInt("page_size", 20)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return
}
