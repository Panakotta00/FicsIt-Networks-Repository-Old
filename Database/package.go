package Database

import (
	"database/sql"
	"net/url"
	"time"
)

type Package struct {
	ID          uint64
	Name        string
	DisplayName string
	Description string
	SourceLink  url.URL
	CreatorID   uint64
	Verified    bool
}

type PackageChange struct {
	ID          uint64
	Name        *string
	DisplayName *string
	Description *string
	SourceLink  *url.URL
	CreatorDate *time.Time
}

type PackageTags struct {
	RepositoryID uint64
	TagID        uint64
}

func PackageGet(db *sql.DB, packageID uint64) (*Package, error) {
	pack := new(Package)
	err := db.QueryRow(
		`SELECT package_id, package_name, package_displayname, package_description, package_creator_id, package_sourcelink, package_verified
			FROM Package WHERE ID=$1`, packageID,
	).Scan(pack.ID, pack.Name, pack.DisplayName, pack.Description, pack.CreatorID, pack.SourceLink, pack.Verified)
	if err != nil {
		return nil, err
	}
	return pack, nil
}
