package Convert

//go:generate go run github.com/jmattheis/goverter/cmd/goverter FINRepository/Convert

import (
	"FINRepository/Database"
	"FINRepository/graph/graphtypes"
	"FINRepository/graph/model"
)

// goverter:converter
// goverter:extend UserIdToUser
// goverter:extend PackageIdToPackage
type ConverterDB interface {
	// goverter:ignore Releases
	// goverter:map CreatorID Creator
	ConvertPackage(Database.Package) model.Package
	ConvertPackageP(*Database.Package) *model.Package

	// goverter:ignore Packages
	// goverter:map EMail Email
	ConvertUser(Database.User) model.User

	ConvertTag(Database.Tag) model.Tag
	ConvertTagP(*Database.Tag) *model.Tag
	ConvertTagA([]Database.Tag) []model.Tag
	ConvertTagPA([]*Database.Tag) []*model.Tag

	// goverter:map PackageID Package
	// goverter:map FINVersion FinVersion
	ConvertRelease(Database.Release) model.Release
	ConvertReleaseP(*Database.Release) *model.Release
	ConvertReleaseA([]Database.Release) []model.Release
}

func UserIdToUser(id Database.ID) model.User {
	return model.User{ID: graphtypes.ID(id)}
}

func PackageIdToPackage(id Database.ID) model.Package {
	return model.Package{ID: graphtypes.ID(id)}
}

// goverter:converter
// goverter:extend UserToUserId
type ConverterModel interface {
	// goverter:map Creator CreatorID
	// goverter:ignore CreatorS
	// goverter:ignore Tags
	ConvertPackage(model.Package) Database.Package
	ConvertPackageP(*model.Package) *Database.Package

	// goverter:ignore ID
	// goverter:ignore Package
	// goverter:ignore Hash
	// goverter:ignore Verified
	// goverter:map FinVersion FINVersion
	ConvertNewRelease(model.NewRelease) Database.Release
}

func UserToUserId(creator *model.User) Database.ID {
	if creator == nil {
		return 0
	}
	return Database.ID(creator.ID)
}
