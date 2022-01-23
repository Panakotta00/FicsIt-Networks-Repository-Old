package UtilReflection

import (
	"FINRepository/Database"
	"reflect"
	"strings"
)

// TODO: Use cache lookup tables for field search
func FindFieldWithMeta(obj interface{}, metaKey string, metaContents string) (*reflect.Value, *reflect.StructField){
	objType := reflect.TypeOf(obj).Elem()
	objValue := reflect.ValueOf(obj).Elem()
	for i := 0; i < objType.NumField(); i++ {
		f := objType.Field(i)
		v := objValue.Field(i)
		meta := f.Tag.Get(metaKey)
		if strings.Contains(meta, metaContents) {
			return &v, &f
		}
	}
	return nil, nil
}

func FindPrimaryKey(obj interface{}) Database.ID {
	v, _ := FindFieldWithMeta(obj, "gorm", "primaryKey")
	if v == nil {
		return 0
	}
	return  Database.ID(v.Int())
}