package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo/v4"
	"log"
	"main/Database"
	"net/http"
	"os"
	"strconv"
)

var db *pgx.Conn

func getPackage(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
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
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid package-id format")
	}
	tags, err := Database.PackageTags(db, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "package not found")
	}
	return c.JSON(http.StatusOK, *tags)
}

func getTag(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
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
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid release-id format")
	}
	release, err := Database.ReleaseGet(db, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "release not found")
	}
	return c.JSON(http.StatusOK, release)
}

func getUser(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user-id format")
	}
	user, err := Database.UserGet(db, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	return c.JSON(http.StatusOK, user)
}

func main() {
	port, err := strconv.Atoi(os.Getenv("FINREPO_DB_PORT"))
	if err != nil || port < 0 {
		log.Fatal("Invalid Port: %v", err)
	}
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", os.Getenv("FINREPO_DB_HOST"), port, os.Getenv("FINREPO_DB_USER"), os.Getenv("FINREPO_DB_PASSWORD"), os.Getenv("FINREPO_DB_DATABASE"))
	db, err = pgx.Connect(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	e := echo.New()

	e.GET("/package/:id", getPackage)
	e.GET("/package/:id/tags", getPackageTags)
	e.GET("/release/:id", getRelease)
	e.GET("/tag/:id", getTag)
	e.GET("/user/:id", getUser)

	e.Logger.Fatal(e.Start(":8000"))
}
