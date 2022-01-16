package Database

import (
	"context"
	"github.com/jackc/pgx/v4"
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

func TagGet(db *pgx.Conn, packageID uint64) (*Tag, error) {
	tag := new(Tag)
	err := db.QueryRow(context.Background(),
		`SELECT tag_id, tag_name, tag_description, tag_verified
			FROM "Repository"."Tag" WHERE tag_id=$1`, packageID,
	).Scan(&tag.ID, &tag.Name, &tag.Description, &tag.Verified)
	if err != nil {
		return nil, err
	}
	return tag, nil
}
