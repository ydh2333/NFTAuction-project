package utils

import "gorm.io/gorm"

// 分页参数
type PageParams struct {
	Page int // 页码（从1开始）
	Size int // 每页条数
}

// 封装分页范围
func Paginate(params PageParams) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		offset := (params.Page - 1) * params.Size
		return tx.Offset(offset).Limit(params.Size)
	}
}
