package helpers

import (
	"manara/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPaginationParams(c *gin.Context) models.PaginationParams {
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}

	perPage, _ := strconv.Atoi(c.Query("per_page"))
	if perPage < 1 {
		perPage = 10
	}
	if perPage > 100 {
		perPage = 100
	}

	return models.PaginationParams{
		Page:    page,
		PerPage: perPage,
	}
}

func Paginate(query *gorm.DB, params models.PaginationParams, result interface{}) (models.Pagination, error) {
	var totalItems int64

	if err := query.Count(&totalItems).Error; err != nil {
		return models.Pagination{}, err
	}

	offset := (params.Page - 1) * params.PerPage
	if err := query.Limit(params.PerPage).Offset(offset).Find(result).Error; err != nil {
		return models.Pagination{}, err
	}

	totalPages := int(totalItems) / params.PerPage
	if int(totalItems)%params.PerPage != 0 {
		totalPages++
	}

	pagination := models.Pagination{
		Page:       params.Page,
		PerPage:    params.PerPage,
		TotalPages: totalPages,
		TotalItems: int(totalItems),
	}

	return pagination, nil
}
