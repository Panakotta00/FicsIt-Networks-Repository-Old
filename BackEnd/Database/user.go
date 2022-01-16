package Database

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type User struct {
	ID          uint64  `json:"id" gorm:";column:user_id;not null;primaryKey"`
	Name        string  `json:"name" gorm:"column:user_name;not null"`
	Bio         string  `json:"bio" gorm:"column:user_bio;not null"`
	GoogleToken *string `json:"googleToken" gorm:"column:user_google_token"`
	Admin       bool    `json:"admin" gorm:"column:user_admin;not nullM;default:false"`
	EMail       string  `json:"email" gorm:"column:user_email;not null"`
	Verified    bool    `json:"verified" gorm:"column:user_verified;not null;default:false"`
}

func (User) TableName() string {
	return "Repository.User"
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
