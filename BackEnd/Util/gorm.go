package Util

import (
	"context"
	"gorm.io/gorm"
)

type ContextDB struct{}

func Paginate(page int, count int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := page * count
		return db.Offset(offset).Limit(count)
	}
}

func ContextWithDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, ContextDB{}, db)
}

func DBFromContext(ctx context.Context) *gorm.DB {
	return ctx.Value(ContextDB{}).(*gorm.DB)
}