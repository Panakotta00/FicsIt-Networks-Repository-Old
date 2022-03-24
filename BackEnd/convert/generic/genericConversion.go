package generic

import (
	"FINRepository/convert/generated"
	database "FINRepository/database"
	"FINRepository/graph/model"
)

func ConvertToDatabase(obj interface{}) interface{} {
	conv := generated.ConverterModelImpl{}
	switch obj := obj.(type) {
	case *model.Package:
		return conv.ConvertPackageP(obj)
	}
	return nil
}

func ConvertToModel(obj interface{}) interface{} {
	conv := generated.ConverterDBImpl{}
	switch obj := obj.(type) {
	case *database.Package:
		return conv.ConvertPackageP(obj)
	}
	return nil
}
