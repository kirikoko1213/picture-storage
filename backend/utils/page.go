package utils

import "picture_storage/model"

func GetPage(page int, pageSize int) model.Pagination {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return model.Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}
