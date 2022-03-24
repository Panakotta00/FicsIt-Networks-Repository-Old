package util

import (
	"context"
	"gorm.io/gorm"
	"reflect"
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

func FilterUpdateFields(fields map[string]interface{}) (newFields map[string]interface{}) {
	newFields = make(map[string]interface{}, len(fields))
	for k, v := range fields {
		if !reflect.ValueOf(v).IsNil() {
			newFields[k] = v
		}
	}
	return
}
