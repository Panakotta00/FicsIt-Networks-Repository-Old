package main

import (
	"FINRepository/Convert/generic"
	"FINRepository/Database"
	"FINRepository/Util"
	UtilReflection "FINRepository/Util/Reflection"
	"FINRepository/auth"
	"FINRepository/dataloader"
	"FINRepository/graph"
	"FINRepository/graph/generated"
	"FINRepository/graph/model"
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
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
	"reflect"
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
	var pack *Database.Package
	if err == nil {
		pack, err = Database.GetPackageByID(db, id)
	} else {
		pack, err = Database.GetPackageByName(db, c.Param("id"))
	}
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
	tags, err := Database.GetPackageTags(db, id)
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
			newCtx := Util.ContextWithDB(ctx.Request().Context(), db)
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

	gqlConfig := generated.Config{Resolvers: &graph.Resolver{}}
	gqlConfig.Directives.IsAdmin = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
		user := ctx.Value("auth").(*Database.User)

		if user == nil || !user.Admin {
			return nil, fmt.Errorf("Access denied")
		}

		return next(ctx)
	}
	gqlConfig.Directives.OwnsOrIsAdmin = func(ctx context.Context, obj interface{}, next graphql.Resolver, owningField string) (interface{}, error) {
		user := ctx.Value("auth").(*Database.User)

		if user == nil {
			return nil, fmt.Errorf("Access denied")
		}

		if user.Admin {
			return next(ctx)
		}

		if reflect.TypeOf(obj) == reflect.TypeOf(&model.User{}) {
			if (Database.ID(obj.(*model.User).ID) == user.ID) {
				return next(ctx)
			} else {
				return nil, fmt.Errorf("Access denied")
			}
		}

		db := Util.DBFromContext(ctx)

		// TODO: Use at boot generated LookUp-Tables instead of direct field search for json
		owningFields := strings.Split(owningField, ".")
		currentObj := obj
		for fieldIndex := 0; fieldIndex < len(owningFields); fieldIndex++ {
			field := owningFields[fieldIndex]
			val := reflect.ValueOf(currentObj).Elem()
			valT := reflect.TypeOf(currentObj).Elem()
			if !val.IsValid() {
				return nil, fmt.Errorf("Access denied")
			}
			for i := 0; i < val.NumField(); i++ {
				f := valT.Field(i)
				v := val.Field(i)
				if f.Tag.Get("json") == field {
					switch v.Kind() {
					case reflect.Int64:
						if int64(user.ID) != v.Int() {
							return nil, fmt.Errorf("Access denied")
						} else {
							return next(ctx)
						}
					case reflect.Ptr:
						if v.IsNil() {
							var dbObj = generic.ConvertToDatabase(currentObj)
							if err := db.Find(&dbObj, UtilReflection.FindPrimaryKey(dbObj)).Error; err != nil {
								return nil, fmt.Errorf("Unable to authorize")
							}
							currentObj = generic.ConvertToModel(dbObj)
							fieldIndex--
						} else {
							currentObj = v.Interface()
						}
					}
					break
				}
			}
		}

		return next(ctx)
	}

	e.Use(auth.AuthenticationMiddleware)

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
	e.Any("/query", echo.WrapHandler(srv))

	e.Logger.Fatal(e.Start(":8000"))
}
