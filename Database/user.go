package Database

type User struct {
	ID          uint64
	Name        string
	Bio         string
	Verified    bool
	GoogleToken *string
	Admin       bool
	EMail       string
}

type UserChange struct {
	ID   uint64
	Name *string
	Bio  *string
}
