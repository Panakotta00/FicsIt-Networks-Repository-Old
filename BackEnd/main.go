package main

import (
	"FINRepository/auth"
	"FINRepository/auth/perm"
	"FINRepository/database"
	"FINRepository/database/cache"
	"FINRepository/dataloader"
	"FINRepository/graph"
	"FINRepository/graph/generated"
	"FINRepository/util"
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/bwmarrin/snowflake"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var db *gorm.DB

func listPackages(c echo.Context) error {
	page, count := util.GetPagination(c)
	packages, err := database.ListPackages(db, page, count)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to list packages")
	}
	return c.JSON(http.StatusOK, packages)
}

func getPackage(c echo.Context) error {
	id, err := util.GetSnowflake(c, "id")
	var pack *database.Package
	if err == nil {
		pack, err = database.GetPackageByID(db, id)
	} else {
		pack, err = database.GetPackageByName(db, c.Param("id"))
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "package not found")
	}
	return c.JSON(http.StatusOK, pack)
}

func getPackageTags(c echo.Context) error {
	id, err := util.GetSnowflake(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid package-id format")
	}
	tags, err := database.GetPackageTags(db, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "package not found")
	}
	return c.JSON(http.StatusOK, *tags)
}

func listPackageReleases(c echo.Context) error {
	id, err := util.GetSnowflake(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid package-id format")
	}
	page, count := util.GetPagination(c)
	releases, err := database.ListPackageReleases(db, id, page, count)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "package not found")
	}
	return c.JSON(http.StatusOK, *releases)
}

func getTags(c echo.Context) error {
	tags, err := database.GetTags(db)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to get all tags")
	}
	return c.JSON(http.StatusOK, tags)
}

func getTag(c echo.Context) error {
	id, err := util.GetSnowflake(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid tag-id format")
	}
	tag, err := database.TagGet(db, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "tag not found")
	}
	return c.JSON(http.StatusOK, tag)
}

func getRelease(c echo.Context) error {
	id, err := util.GetSnowflake(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid release-id format")
	}
	release, err := database.ReleaseGet(db, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "release not found")
	}
	return c.JSON(http.StatusOK, release)
}

func listUsers(c echo.Context) error {
	page, count := util.GetPagination(c)
	users, err := database.ListUsers(db, page, count)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to list users")
	}
	return c.JSON(http.StatusOK, *users)
}

func getUser(c echo.Context) error {
	id, err := util.GetSnowflake(c, "id")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user-id format")
	}
	user, err := database.UserGet(db, id)
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

	err = db.SetupJoinTable(&database.Package{}, "Tags", &database.PackageTag{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&database.User{}, &database.UserChange{}, &database.Package{}, &database.PackageChange{}, &database.Tag{}, &database.Release{}, &database.ReleaseChange{})
	if err != nil {
		log.Fatal(err)
	}

	port, err = strconv.Atoi(os.Getenv("FINREPO_SPICEDB_PORT"))
	if err != nil || port < 0 {
		log.Fatal("Invalid SpiceDB Port: %v", err)
	}
	connectionString = fmt.Sprintf("%s:%d", os.Getenv("FINREPO_SPICEDB_HOST"), port)

	spicedb, err := auth.NewSpiceDBAuthorizer(connectionString, os.Getenv("FINREPO_SPICEDB_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to start SpiceDB Client: %v", err)
	}

	node, err := strconv.ParseInt(os.Getenv("FINREPO_NODE"), 10, 64)
	if err != nil || node < 0 {
		log.Fatal("Invalid Node: %v", err)
	}
	idGen, err := snowflake.NewNode(node)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			newCtx := perm.CtxWithAuthorizer(ctx.Request().Context(), spicedb)
			ctx.SetRequest(ctx.Request().WithContext(newCtx))
			return next(ctx)
		}
	})

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			newCtx := util.ContextWithDB(ctx.Request().Context(), db)
			ctx.SetRequest(ctx.Request().WithContext(newCtx))
			return next(ctx)
		}
	})

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			newCtx := context.WithValue(ctx.Request().Context(), "snowflake", idGen)
			ctx.SetRequest(ctx.Request().WithContext(newCtx))
			return next(ctx)
		}
	})

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			newCtx := cache.CtxWithDBCache(ctx.Request().Context())
			ctx.SetRequest(ctx.Request().WithContext(newCtx))
			return next(ctx)
		}
	})

	gqlConfig := generated.Config{Resolvers: &graph.Resolver{}}
	gqlConfig.Directives.IsAdmin = graph.IsAdminDirective
	gqlConfig.Directives.OwnsOrIsAdmin = graph.OwnsOrIsAdminDirective
	gqlConfig.Directives.Permission = graph.PermissionDirective

	e.Use(dataloader.Middleware)

	e.GET("/package", listPackages)
	e.GET("/package/:id", getPackage)
	e.GET("/package/:id/tags", getPackageTags)
	e.GET("/package/:id/releases", listPackageReleases)
	e.GET("/release/:id", getRelease)
	e.GET("/tag", getTags)
	e.GET("/tag/:id", getTag)
	e.GET("/user", listUsers)
	e.GET("/user/:id", getUser)
	e.POST("/oauth", auth.OAuth2Request)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(gqlConfig))

	e.Any("/playground", echo.WrapHandler(playground.Handler("GraphQL playground", "/query")))
	e.Any("/query", echo.WrapHandler(srv), auth.AuthenticationMiddleware)

	e.Logger.Fatal(e.Start(":8000"))
}
