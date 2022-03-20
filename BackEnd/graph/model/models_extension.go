package model

import "strconv"

func (r *Release) GetType() string {
	return "release"
}

func (r *Release) GetID() string {
	return strconv.FormatInt(int64(r.ID), 10)
}

func (r *Release) IsVerified() bool {
	return r.Verified
}

func (u *User) GetType() string {
	return "user"
}

func (u *User) GetID() string {
	return strconv.FormatInt(int64(u.ID), 10)
}

func (u *User) IsVerified() bool {
	return u.Verified
}
