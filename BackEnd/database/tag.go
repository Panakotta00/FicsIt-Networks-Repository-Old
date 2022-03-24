package database

import (
	"FINRepository/auth/perm"
	"FINRepository/util"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type Tag struct {
	ID          ID         `json:"id" gorm:"column:tag_id;not null;primaryKey"`
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

func TagGet(db *gorm.DB, tagID int64) (*Tag, error) {
	tag := new(Tag)
	if err := db.First(tag, tagID).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

func CreateTag(ctx context.Context, tagName string, tagDescription string) (*Tag, error) {
	success, _ := perm.AuthorizerFromCtx(ctx).Authorize(ctx, &perm.Global, UserFromCtx(ctx), "edit_tags")
	if !success {
		return nil, fmt.Errorf("access denied")
	}
	tag := Tag{ID: ID(util.GetSnowflakeFromCTX(ctx).Generate().Int64()), Name: tagName, Description: tagDescription}
	err := util.DBFromContext(ctx).Create(&tag).Error
	return &tag, err
}

func UpdateTag(ctx context.Context, tagId ID, tagName *string, tagDescription *string) (bool, error) {
	success, _ := perm.AuthorizerFromCtx(ctx).Authorize(ctx, &perm.Global, UserFromCtx(ctx), "edit_tags")
	if !success {
		return false, fmt.Errorf("access denied")
	}
	fields := util.FilterUpdateFields(map[string]interface{}{"tag_name": tagName, "tag_description": tagDescription})
	result := util.DBFromContext(ctx).Model(&Tag{}).Where("tag_id = ?", tagId).Updates(fields)
	return result.Error == nil && result.RowsAffected == 1, nil
}

func DeleteTag(ctx context.Context, tagId ID) (bool, error) {
	success, _ := perm.AuthorizerFromCtx(ctx).Authorize(ctx, &perm.Global, UserFromCtx(ctx), "edit_tags")
	if !success {
		return false, fmt.Errorf("access denied")
	}
	util.DBFromContext(ctx).Where("tag_id = ?", tagId).Delete(&PackageTag{})
	result := util.DBFromContext(ctx).Where("tag_id = ?", tagId).Delete(&Tag{})
	return result.Error == nil && result.RowsAffected == 1, nil
}
