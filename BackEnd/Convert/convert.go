package Convert

//go:generate go run github.com/jmattheis/goverter/cmd/goverter FINRepository/Convert

import (
	"FINRepository/Database"
	"FINRepository/graph/model"
)

// goverter:converter
// goverter:extend UserIdToUser
type Converter interface {
	// goverter:ignore Tags
	// goverter:ignore Releases
	// goverter:map CreatorID Creator
	ConvertPackage(Database.Package) model.Package

	// goverter:ignore Packages
	// goverter:map EMail Email
	ConvertUser(Database.User) model.User
}

func UserIdToUser(id uint64) model.User {
	return model.User{ID: id}
}
