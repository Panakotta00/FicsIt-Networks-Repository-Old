package database

import (
	"FINRepository/auth/perm"
	"FINRepository/util"
	"context"
	"gorm.io/gorm"
	"strconv"
)

type Release struct {
	ID          ID       `json:"id" gorm:"column:release_id;not null;primaryKey"`
	PackageID   ID       `json:"packageId" gorm:"column:package_id;not null"`
	Name        string   `json:"name" gorm:"column:release_name;not null"`
	Description string   `json:"description" gorm:"column:release_description;not null"`
	SourceLink  string   `json:"sourceLink" gorm:"column:release_source_link;not null"`
	Version     string   `json:"version" gorm:"column:release_version;not null"`
	FINVersion  string   `json:"finVersion" gorm:"column:release_fin_version;not null"`
	Verified    bool     `json:"verified" gorm:"column:release_verified;not null;default:false"`
	Hash        string   `json:"hash" gorm:"column:release_hash;not null"`
	Package     *Package `json:"package,omitempty" gorm:"foreignKey:PackageID"`
	ZedToken    string   `gorm:"column:release_zedtoken"`
}

func (r *Release) GetType() string {
	return "release"
}

func (r *Release) GetID() string {
	return strconv.FormatInt(int64(r.ID), 10)
}

type ReleaseChange struct {
	ID          ID       `json:"id" gorm:"column:release_change_id;not null;primaryKey"`
	Name        *string  `json:"name" gorm:"column:release_name"`
	Description *string  `json:"description" gorm:"column:release_description"`
	Release     *Release `json:"release,omitempty" gorm:"foreignKey:release_change_id"`
}

func ReleaseGet(db *gorm.DB, releaseId int64) (*Release, error) {
	release := new(Release)
	if err := db.First(&release, releaseId).Error; err != nil {
		return nil, err
	}
	return release, nil
}

func CreateRelease(ctx context.Context, packageId ID, name string, description string, sourceURL string, version string, finVersion string) (*Release, error) {
	id := util.GetSnowflakeFromCTX(ctx).Generate().Int64()
	authorizer := perm.AuthorizerFromCtx(ctx)
	release := Release{
		ID:          ID(id),
		PackageID:   packageId,
		Name:        name,
		Description: description,
		SourceLink:  sourceURL,
		Version:     version,
		FINVersion:  finVersion,
	}
	token, err := authorizer.AddRelation(ctx, &release, &perm.AuthorizableGeneric{"package", strconv.FormatInt(int64(packageId), 10)}, "package")
	if err != nil {
		return nil, err
	}
	release.ZedToken = token
	err = util.DBFromContext(ctx).Create(&release).Error
	return &release, err
}
