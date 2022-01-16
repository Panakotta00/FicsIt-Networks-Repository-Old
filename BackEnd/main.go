package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"main/Database"
	"main/Util"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var db *gorm.DB

func listPackages(c echo.Context) error {
	page, count := Util.GetPagination(c)
	packages, err := Database.ListPackages(db, page, count)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to list packages")
	}
	return c.JSON(http.StatusOK, packages)
}

func getPackage(c echo.Context) error {
	id, err := Util.GetSnowflake(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid package-id format")
	}
	pack, err := Database.PackageGet(db, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "package not found")
	}
	return c.JSON(http.StatusOK, pack)
}

func getPackageTags(c echo.Context) error {
	id, err := Util.GetSnowflake(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid package-id format")
	}
	tags, err := Database.PackageTags(db, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "package not found")
	}
	return c.JSON(http.StatusOK, *tags)
}

func listPackageReleases(c echo.Context) error {
	id, err := Util.GetSnowflake(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid package-id format")
	}
	page, count := Util.GetPagination(c)
	releases, err := Database.ListPackageReleases(db, id, page, count)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "package not found")
	}
	return c.JSON(http.StatusOK, *releases)
}

func getTags(c echo.Context) error {
	tags, err := Database.GetTags(db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to get all tags")
	}
	return c.JSON(http.StatusOK, tags)
}

func getTag(c echo.Context) error {
	id, err := Util.GetSnowflake(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid tag-id format")
	}
	tag, err := Database.TagGet(db, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "tag not found")
	}
	return c.JSON(http.StatusOK, tag)
}

func getRelease(c echo.Context) error {
	id, err := Util.GetSnowflake(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid release-id format")
	}
	release, err := Database.ReleaseGet(db, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "release not found")
	}
	return c.JSON(http.StatusOK, release)
}

func listUsers(c echo.Context) error {
	page, count := Util.GetPagination(c)
	users, err := Database.ListUsers(db, page, count)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to list users")
	}
	return c.JSON(http.StatusOK, *users)
}

func getUser(c echo.Context) error {
	id, err := Util.GetSnowflake(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user-id format")
	}
	user, err := Database.UserGet(db, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	return c.JSON(http.StatusOK, user)
}

type CustomNamer struct {
	schema.NamingStrategy
	TablePrefix   string
	SingularTable bool
}

func (n CustomNamer) TableName(table string) string {
	return n.TablePrefix + table
}

func (n CustomNamer) SchemaName(table string) string {
	return strings.TrimPrefix(table, n.TablePrefix)
}

func (_ CustomNamer) ColumnName(table, column string) string {
	return strings.ToLower(column)
}

func (_ CustomNamer) JoinTableName(table string) string {
	return "Repository." + table
}

func (_ CustomNamer) RelationshipFKName(rel schema.Relationship) string {
	return "fk_" + rel.Schema.Name + "_" + rel.Name
}

func (_ CustomNamer) CheckerName(table, column string) string {
	return column
}

func (_ CustomNamer) IndexName(table, column string) string {
	return column
}

func main() {
	port, err := strconv.Atoi(os.Getenv("FINREPO_DB_PORT"))
	if err != nil || port < 0 {
		log.Fatal("Invalid Port: %v", err)
	}
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", os.Getenv("FINREPO_DB_HOST"), port, os.Getenv("FINREPO_DB_USER"), os.Getenv("FINREPO_DB_PASSWORD"), os.Getenv("FINREPO_DB_DATABASE"))

	db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{
		NamingStrategy: CustomNamer{
			TablePrefix:   "Repository.",
			SingularTable: false,
		},
	})
	if err != nil {
		log.Fatal("failed to connect database: %v", err)
	}

	err = db.SetupJoinTable(&Database.Package{}, "Tags", &Database.PackageTag{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&Database.User{}, &Database.UserChange{}, &Database.Package{}, &Database.PackageChange{}, &Database.Tag{}, &Database.Release{}, &Database.ReleaseChange{})
	if err != nil {
		log.Fatal(err)
	}

	tags, err := Database.PackageTags(db, 123)
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.Marshal(tags)
	log.Println(string(data))

	e := echo.New()

	e.GET("/package", listPackages)
	e.GET("/package/:id", getPackage)
	e.GET("/package/:id/tags", getPackageTags)
	e.GET("/package/:id/releases", listPackageReleases)
	e.GET("/release/:id", getRelease)
	e.GET("/tag", getTags)
	e.GET("/tag/:id", getTag)
	e.GET("/user", listUsers)
	e.GET("/user/:id", getUser)

	e.Logger.Fatal(e.Start(":8000"))
}
