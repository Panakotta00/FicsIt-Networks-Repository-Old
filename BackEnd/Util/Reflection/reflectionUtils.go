package UtilReflection

import (
	"FINRepository/Database"
	"reflect"
	"strings"
)

type fieldToMetaSubKey struct {
	Type reflect.Type
	MetaKey string
	MetaSubstring string
}
var fieldToMetaSubMap = map[fieldToMetaSubKey]int{}
func FindFieldWithMetaSubstring(obj interface{}, metaKey string, substring string) (*reflect.Value, *reflect.StructField){
	objType := reflect.TypeOf(obj).Elem()
	objValue := reflect.ValueOf(obj).Elem()
	key := fieldToMetaSubKey{objType, metaKey, substring}
	fieldIndex, ok := fieldToMetaSubMap[key]

	if !ok {
		for i := 0; i < objType.NumField(); i++ {
			f := objType.Field(i)
			meta := f.Tag.Get(metaKey)
			if strings.Contains(meta, substring) {
				fieldToMetaSubMap[key] = i
				break
			}
		}
	}

	if fieldIndex < 0 {
		return nil, nil
	} else {
		v := objValue.Field(fieldIndex)
		f := objType.Field(fieldIndex)
		return &v, &f
	}
}

type fieldToMetaKey struct {
	Type reflect.Type
	MetaKey string
	MetaContents string
}
var fieldToMetaMap = map[fieldToMetaKey]int{}
func FindFieldWithMeta(obj interface{}, metaKey string, contents string) (*reflect.Value, *reflect.StructField) {
	objType := reflect.TypeOf(obj).Elem()
	objValue := reflect.ValueOf(obj).Elem()
	key := fieldToMetaKey{objType, metaKey, contents}
	fieldIndex, ok := fieldToMetaMap[key]

	if !ok {
		for i := 0; i < objType.NumField(); i++ {
			f := objType.Field(i)
			meta := f.Tag.Get(metaKey)
			if meta == contents {
				fieldToMetaMap[key] = i
				break
			}
		}
	}

	if fieldIndex < 0 {
		return nil, nil
	} else {
		v := objValue.Field(fieldIndex)
		f := objType.Field(fieldIndex)
		return &v, &f
	}
}

func FindPrimaryKey(obj interface{}) Database.ID {
	v, _ := FindFieldWithMetaSubstring(obj, "gorm", "primaryKey")
	if v == nil {
		return 0
	}
	return  Database.ID(v.Int())
}