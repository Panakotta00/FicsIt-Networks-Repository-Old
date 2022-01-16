package Database

import (
	"gorm.io/gorm"
	"main/Util"
)

type PackageLimited struct {
	ID          uint64  `json:"id" gorm:"column:package_id;not null;primaryKey"`
	Name        string  `json:"name" gorm:"column:package_name;not null;unique"`
	DisplayName string  `json:"displayName" gorm:"column:package_displayname;not null"`
	Description string  `json:"description" gorm:"column:package_description;not null"`
	SourceLink  *string `json:"sourceLink" gorm:"column:package_sourcelink;"`
	CreatorID   uint64  `json:"creatorId" gorm:"column:package_creator_id;not null"`
	CreatorS    *User   `json:"creator,omitempty" gorm:"foreignKey:CreatorID"`
	Tags        []*Tag  `json:"tags,omitempty" gorm:"many2many:Package_Tag;foreignKey:package_id;joinForeignKey:package_id;References:tag_id;joinReferences:tag_id"`
}

func (PackageLimited) TableName() string {
	return "Repository.Package"
}

type Package struct {
	PackageLimited
	Verified bool `json:"verified" gorm:"column:package_verified;not null;default:false"`
}

type PackageChange struct {
	ID          uint64   `json:"id" gorm:"column:package_change_id;not null;primaryKey"`
	Name        *string  `json:"name" gorm:"column:package_name"`
	DisplayName *string  `json:"displayName" gorm:"column:package_displayname"`
	Description *string  `json:"description" gorm:"column:package_description"`
	SourceLink  *string  `json:"sourceLink" gorm:"column:package_sourcelink"`
	Package     *Package `json:"package,omitempty" gorm:"foreignKey:package_change_id"`
}

func (PackageChange) TableName() string {
	return "Repository.Package_Change"
}

type PackageTag struct {
	PackageID uint64 `json:"packageId" gorm:"column:package_id;primaryKey"`
	TagID     uint64 `json:"tagId" gorm:"column:tag_id;primaryKey"`
}

func (PackageTag) TableName() string {
	return "Repository.Package_Tag"
}

func ListPackages(db *gorm.DB, page int, count int) (*[]*Package, error) {
	var packages *[]*Package
	if err := db.Scopes(Util.Paginate(page, count)).Find(&packages).Error; err != nil {
		return nil, err
	}
	return packages, nil
}

func PackageGet(db *gorm.DB, packageID uint64) (*Package, error) {
	pack := new(Package)
	if err := db.First(&pack, packageID).Error; err != nil {
		return nil, err
	}
	return pack, nil
}

func PackageTags(db *gorm.DB, packageId uint64) (*[]*Tag, error) {
	var pack Package
	if err := db.Preload("Tags").Select("ID").First(&pack, packageId).Error; err != nil {
		return nil, err
	}
	return &pack.Tags, nil
}
