package Database

import (
	"gorm.io/gorm"
)

type Release struct {
	ID          uint64   `json:"id" gorm:"column:release_id;not null;primaryKey"`
	PackageID   uint64   `json:"packageId" gorm:"column:package_id;not null"`
	Name        string   `json:"name" gorm:"column:release_name;not null"`
	Description string   `json:"description" gorm:"column:release_description;not null"`
	SourceLink  string   `json:"sourceLink" gorm:"column:release_source_link;not null"`
	Version     string   `json:"version" gorm:"column:release_version;not null"`
	FINVersion  string   `json:"finVersion" gorm:"column:release_fin_version;not null"`
	Verified    bool     `json:"verified" gorm:"column:release_verified;not null;default:false"`
	Hash        string   `json:"hash" gorm:"column:release_hash;not null"`
	Package     *Package `json:"package,omitempty" gorm:"foreignKey:PackageID"`
}

type ReleaseChange struct {
	ID          uint64  `json:"id" gorm:"column:release_change_id;not null;primaryKey"`
	Name        *string `json:"name" gorm:"column:release_name"`
	Description *string `json:"description" gorm:"column:release_description"`
	Release     *Release `json:"release,omitempty" gorm:"foreignKey:release_change_id"`
}

func ReleaseGet(db *gorm.DB, releaseId uint64) (*Release, error) {
	release := new(Release)
	if err := db.First(&release, releaseId).Error; err != nil {
		return nil, err
	}
	return release, nil
}
