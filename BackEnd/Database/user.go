package Database

import (
	"FINRepository/Util"
	"gorm.io/gorm"
)

type User struct {
	ID          uint64  `json:"id" gorm:"column:user_id;not null;primaryKey"`
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
	ID   uint64  `json:"id" gorm:"column:user_change_id;not null;primaryKey"`
	Name *string `json:"name" gorm:"column:user_name"`
	Bio  *string `json:"bio" gorm:"column:user_bio"`
	User *User   `json:"user,omitempty" gorm:"foreignKey:user_change_id"`
}

func (UserChange) TableName() string {
	return "Repository.User_Change"
}

func ListUsers(db *gorm.DB, page int, count int) (*[]*User, error) {
	var users *[]*User
	if err := db.Scopes(Util.Paginate(page, count)).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func UserGet(db *gorm.DB, userId uint64) (*User, error) {
	user := new(User)
	if err := db.First(user, userId).Error; err != nil {
		return nil, err
	}
	return user, nil
}
