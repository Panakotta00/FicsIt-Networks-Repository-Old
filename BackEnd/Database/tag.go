package Database

import (
	"gorm.io/gorm"
)

type Tag struct {
	ID          uint64     `json:"id" gorm:"column:tag_id;not null;primaryKey"`
	Name        string     `json:"name" gorm:"column:tag_name;not null"`
	Description string     `json:"description" gorm:"column:tag_description;not null"`
	Verified    bool       `json:"verified" gorm:"column:tag_verified;not null;default:false"`
	Packages    []*Package `json:"packages,omitempty" gorm:"many2many:Package_Tag;foreignKey:tag_id;joinForeignKey:tag_id;References:package_id;joinReferences:package_id"`
}

func (Tag) TableName() string {
	return "Repository.Tag"
}

func GetTags(db *gorm.DB) (*[]*Tag, error) {
	var tags *[]*Tag
	if err := db.Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func TagGet(db *gorm.DB, tagID uint64) (*Tag, error) {
	tag := new(Tag)
	if err := db.First(tag, tagID).Error; err != nil {
		return nil, err
	}
	return tag, nil
}
