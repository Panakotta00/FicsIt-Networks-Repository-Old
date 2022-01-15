package Database

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type Release struct {
	ID          uint64 `json:"id"`
	Package     uint64 `json:"package"`
	Name        string `json:"name"`
	Description string `json:"description"`
	SourceLink  string `json:"sourceLink"`
	Version     string `json:"version"`
	FINVersion  string `json:"finVersion"`
	Verified    bool   `json:"verified"`
	Hash        string `json:"hash"`
}

type ReleaseChange struct {
	ID          uint64  `json:"id""`
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

func ReleaseGet(db *pgx.Conn, releaseId uint64) (*Release, error) {
	release := new(Release)
	err := db.QueryRow(context.Background(),
		`SELECT release_id, package_id, release_name, release_description, release_source_link, release_version, release_fin_version, release_verified, release_hash
			FROM "Repository"."Release" WHERE release_id=$1`, releaseId,
	).Scan(&release.ID, &release.Package, &release.Name, &release.Description, &release.SourceLink, &release.Version, &release.FINVersion, &release.Verified, &release.Hash)
	if err != nil {
		return nil, err
	}
	return release, nil
}
