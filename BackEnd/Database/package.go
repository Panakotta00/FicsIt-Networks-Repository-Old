package Database

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type PackageLimited struct {
	ID          uint64  `json:"id"`
	Name        string  `json:"name"`
	DisplayName string  `json:"displayName"`
	Description string  `json:"description"`
	SourceLink  *string `json:"sourceLink"`
	CreatorID   uint64  `json:"creatorId"`
}

type Package struct {
	PackageLimited
	Verified bool `json:"verified"`
}

type PackageChange struct {
	ID          uint64  `json:"id"`
	Name        *string `json:"name"`
	DisplayName *string `json:"displayName"`
	Description *string `json:"description"`
	SourceLink  *string `json:"sourceLink"`
}

type PackageTag struct {
	RepositoryID uint64 `json:"repositoryId"`
	TagID        uint64 `json:"tagId"`
}

func PackageGet(db *pgx.Conn, packageID uint64) (*Package, error) {
	pack := new(Package)
	err := db.QueryRow(context.Background(),
		`SELECT package_id, package_name, package_displayname, package_description, package_creator_id, package_sourcelink, package_verified
			FROM "Repository"."Package" WHERE package_id=$1`, packageID,
	).Scan(&pack.ID, &pack.Name, &pack.DisplayName, &pack.Description, &pack.CreatorID, &pack.SourceLink, &pack.Verified)
	if err != nil {
		return nil, err
	}
	return pack, nil
}

func PackageTags(db *pgx.Conn, packageId uint64) (*[]Tag, error) {
	rows, err := db.Query(context.Background(), `
		WITH tags AS (
			SELECT tag_id
			FROM "Repository"."Package_Tag"
			WHERE package_id=$1
		)
		SELECT tag_id, tag_name, tag_description, tag_verified
		FROM "Repository"."Tag"
		WHERE "Tag".tag_id IN (SELECT * FROM tags)`, packageId)
	if err != nil {
		return nil, err
	}
	tags := make([]Tag, 0)
	for rows.Next() {
		var tag Tag
		err = rows.Scan(&tag.ID, &tag.Name, &tag.Description, &tag.Verified)
		if err != nil {
			continue
		}
		tags = append(tags, tag)
	}
	return &tags, nil
}
