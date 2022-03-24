package convert

//go:generate go run github.com/jmattheis/goverter/cmd/goverter FINRepository/Convert

import (
	"FINRepository/database"
	"FINRepository/graph/graphtypes"
	"FINRepository/graph/model"
)

// goverter:converter
// goverter:extend UserIdToUser
// goverter:extend PackageIdToPackage
type ConverterDB interface {
	// goverter:ignore Releases
	// goverter:map CreatorID Creator
	ConvertPackage(database.Package) model.Package
	ConvertPackageP(*database.Package) *model.Package

	// goverter:ignore Packages
	// goverter:map EMail Email
	ConvertUser(database.User) model.User
	ConvertUserP(*database.User) *model.User

	ConvertTag(database.Tag) model.Tag
	ConvertTagP(*database.Tag) *model.Tag
	ConvertTagA([]database.Tag) []model.Tag
	ConvertTagPA([]*database.Tag) []*model.Tag

	// goverter:map PackageID Package
	// goverter:map FINVersion FinVersion
	ConvertRelease(database.Release) model.Release
	ConvertReleaseP(*database.Release) *model.Release
	ConvertReleaseA([]database.Release) []model.Release
}

func UserIdToUser(id database.ID) model.User {
	return model.User{ID: graphtypes.ID(id)}
}

func PackageIdToPackage(id database.ID) model.Package {
	return model.Package{ID: graphtypes.ID(id)}
}

// goverter:converter
// goverter:extend UserToUserId
type ConverterModel interface {
	// goverter:map Creator CreatorID
	// goverter:ignore CreatorS
	// goverter:ignore Tags
	ConvertPackage(model.Package) database.Package
	ConvertPackageP(*model.Package) *database.Package

	// goverter:ignore ID
	// goverter:ignore Package
	// goverter:ignore Hash
	// goverter:ignore Verified
	// goverter:map FinVersion FINVersion
	ConvertNewRelease(model.NewRelease) database.Release
}

func UserToUserId(creator *model.User) database.ID {
	if creator == nil {
		return 0
	}
	return database.ID(creator.ID)
}
