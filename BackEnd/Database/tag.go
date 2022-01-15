package Database

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type Tag struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Verified    bool   `json:"verified"`
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