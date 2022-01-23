package graph

import (
	"FINRepository/Convert/generic"
	"FINRepository/Database"
	"FINRepository/Util"
	UtilReflection "FINRepository/Util/Reflection"
	"FINRepository/graph/model"
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"reflect"
	"regexp"
	"strings"
)

var owningFieldNameRegex = regexp.MustCompile("\\[(\\w+),(\\w+)]")

func modelByName(modelName string) interface{} {
	switch modelName {
	case "Package":
		return &model.Package{}
	default:
		return nil
	}
}

func IsAdminDirective(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	user := ctx.Value("auth").(*Database.User)

	if user == nil || !user.Admin {
		return nil, fmt.Errorf("Access denied")
	}

	return next(ctx)
}

func OwnsOrIsAdminDirective(ctx context.Context, obj interface{}, next graphql.Resolver, owningField string) (interface{}, error) {
	user := ctx.Value("auth").(*Database.User)

	if user == nil {
		return nil, fmt.Errorf("Access denied")
	}

	if user.Admin {
		return next(ctx)
	}

	if reflect.TypeOf(obj) == reflect.TypeOf(&model.User{}) {
		if Database.ID(obj.(*model.User).ID) == user.ID {
			return next(ctx)
		} else {
			return nil, fmt.Errorf("Access denied")
		}
	}

	db := Util.DBFromContext(ctx)

	// TODO: Use at boot generated LookUp-Tables instead of direct field search for json
	owningFields := strings.Split(owningField, ".")
	currentObj := obj
	for fieldIndex := 0; fieldIndex < len(owningFields); fieldIndex++ {
		objValue := reflect.ValueOf(currentObj)
		if objValue.Kind() != reflect.Map {
			objValue = objValue.Elem()
		}
		if !objValue.IsValid() {
			return nil, fmt.Errorf("Access denied")
		}

		fieldName := owningFields[fieldIndex]
		var fieldModel *string
		if match := owningFieldNameRegex.FindStringSubmatch(fieldName); match != nil {
			fieldName = match[1]
			fieldModel = &match[2]
		}

		var fieldValue *reflect.Value
		if objValue.Kind() == reflect.Map {
			value := objValue.MapIndex(reflect.ValueOf(fieldName)).Elem()
			fieldValue = &value
		} else {
			fieldValue, _ = UtilReflection.FindFieldWithMeta(currentObj, "json", fieldName)
		}
		if fieldValue == nil {
			return nil, fmt.Errorf("Access denied")
		}

		switch fieldValue.Kind() {
		case reflect.Int64:
			if fieldModel != nil {
				modelTemplate := modelByName(*fieldModel)
				var dbObj = generic.ConvertToDatabase(modelTemplate)
				if err := db.Find(&dbObj, fieldValue.Int()).Error; err != nil {
					return nil, fmt.Errorf("Unable to authorize")
				}
				currentObj = generic.ConvertToModel(dbObj)
			} else if int64(user.ID) != fieldValue.Int() {
				return nil, fmt.Errorf("Access denied")
			} else {
				return next(ctx)
			}
			break
		case reflect.Ptr:
			if fieldValue.IsNil() {
				var dbObj = generic.ConvertToDatabase(currentObj)
				if err := db.Find(&dbObj, UtilReflection.FindPrimaryKey(dbObj)).Error; err != nil {
					return nil, fmt.Errorf("Unable to authorize")
				}
				currentObj = generic.ConvertToModel(dbObj)
				fieldIndex--
			} else {
				currentObj = fieldValue.Interface()
			}
			break
		default:
			return nil, fmt.Errorf("Unable to authorize, Invalid GQL Directive Param, invalid data-type")
		}
	}

	return nil, fmt.Errorf("Unable to authorize, Invalid GQL Directive Param, incomplete")
}
