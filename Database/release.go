package Database

import (
	"github.com/Masterminds/semver"
	"net/url"
)

type Release struct {
	ID          uint64
	Package     uint64
	Name        string
	Description string
	SourceLink  url.URL
	Version     semver.Version
	FINVersion  semver.Version
	Verified    bool
	Hash        uint64
}

type ReleaseChange struct {
	ID          uint64
	Name        *string
	Description *string
}
