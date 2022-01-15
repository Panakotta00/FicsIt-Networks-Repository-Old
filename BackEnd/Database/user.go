package Database

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type User struct {
	ID          uint64  `json:"id"`
	Name        string  `json:"name"`
	Bio         string  `json:"bio"`
	GoogleToken *string `json:"googleToken"`
	Admin       bool    `json:"admin"`
	EMail       string  `json:"email"`
	Verified    bool    `json:"verified"`
}

type UserChange struct {
	ID   uint64
	Name *string
	Bio  *string
}

func UserGet(db *pgx.Conn, userId uint64) (*User, error) {
	user := new(User)
	err := db.QueryRow(context.Background(),
		`SELECT user_id, user_name, user_bio, user_google_token, user_admin, user_email, user_verified
			FROM "Repository"."User" WHERE user_id=$1`, userId,
	).Scan(&user.ID, &user.Name, &user.Bio, &user.GoogleToken, &user.Admin, &user.EMail, &user.Verified)
	if err != nil {
		return nil, err
	}
	return user, nil
}
