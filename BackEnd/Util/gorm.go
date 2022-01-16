package Util

import "gorm.io/gorm"

func Paginate(page int, count int) func(*gorm.DB)*gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := page * count
		return db.Offset(offset).Limit(count)
	}
}
