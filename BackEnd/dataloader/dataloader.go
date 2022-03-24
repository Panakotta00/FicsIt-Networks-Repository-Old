package dataloader

import (
	"FINRepository/convert/generated"
	"FINRepository/database"
	"FINRepository/database/cache"
	"FINRepository/graph/graphtypes"
	"FINRepository/graph/model"
	"FINRepository/util"
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"log"
	"time"
)

//go:generate go run github.com/vektah/dataloaden UserLoader FINRepository/graph/graphtypes.ID *FINRepository/graph/model.User
//go:generate go run github.com/vektah/dataloaden TagsByPackageLoader FINRepository/graph/graphtypes.ID []FINRepository/graph/graphtypes.ID
//go:generate go run github.com/vektah/dataloaden TagLoader FINRepository/graph/graphtypes.ID *FINRepository/graph/model.Tag
//go:generate go run github.com/vektah/dataloaden ReleasesByPackageLoader FINRepository/graph/graphtypes.ID []*FINRepository/graph/model.Release
//go:generate go run github.com/vektah/dataloaden PackageByIdLoader FINRepository/graph/graphtypes.ID *FINRepository/graph/model.Package
//go:generate go run github.com/vektah/dataloaden PackagesByTagLoader FINRepository/graph/graphtypes.ID []FINRepository/graph/graphtypes.ID
//go:generate go run github.com/vektah/dataloaden PackagesByUserLoader FINRepository/graph/graphtypes.ID []*FINRepository/graph/model.Package

type loadersKey struct{}

type Loaders struct {
	UserById          UserLoader
	TagsByPackage     TagsByPackageLoader
	TagById           TagLoader
	ReleasesByPackage ReleasesByPackageLoader
	PackageById       PackageByIdLoader
	PackagesByTag     PackagesByTagLoader
	PackagesByUser    PackagesByUserLoader
}

func Middleware(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		nCtx := context.WithValue(ctx, loadersKey{}, &Loaders{
			UserById: UserLoader{
				maxBatch: 100,
				wait:     time.Millisecond,
				fetch: func(ids []graphtypes.ID) (users []*model.User, errs []error) {
					users = make([]*model.User, len(ids))
					errs = make([]error, len(ids))
					dbCache := cache.DBCacheFromCtx(ctx)

					// try to load users from cache
					var idsToQueryMap = map[graphtypes.ID]int{}
					var idsToQuery []graphtypes.ID
					for i, id := range ids {
						if user := dbCache.GetByPK(&database.User{}, database.ID(id)); user == nil {
							idsToQueryMap[id] = i
							idsToQuery = append(idsToQuery, id)
						} else {
							conv := generated.ConverterDBImpl{}
							users[i] = conv.ConvertUserP((*user).(*database.User))
						}
					}
					if len(idsToQuery) < 1 {
						return
					}

					// query non-cached users from DB
					var dbUsers []*database.User
					if err := util.DBFromContext(ctx).Find(&dbUsers, idsToQuery).Error; err != nil {
						e := errors.New("unable to get users by ids")
						for i := 0; i < len(ids); i++ {
							errs[i] = e
						}
						return
					}

					// cache and convert non-cached users
					var idMap = make(map[graphtypes.ID]model.User, len(dbUsers))
					conv := generated.ConverterDBImpl{}
					for _, user := range dbUsers {
						dbCache.Add(user)
						idMap[graphtypes.ID(user.ID)] = conv.ConvertUser(*user)
					}
					for id, i := range idsToQueryMap {
						user, ok := idMap[id]
						if ok {
							users[i] = &user
						} else {
							errs[i] = errors.New("unable to find user")
						}
					}

					return
				},
			},
			TagsByPackage: TagsByPackageLoader{
				maxBatch: 100,
				wait:     time.Millisecond,
				fetch: func(ids []graphtypes.ID) (packageTags [][]graphtypes.ID, errs []error) {
					errs = make([]error, len(ids))
					packageTags = make([][]graphtypes.ID, len(ids))

					var dbTags []*database.PackageTag
					if err := util.DBFromContext(ctx).Where("package_id IN ?", ids).Find(&dbTags).Error; err != nil {
						log.Printf("Error: %v", err)
						e := errors.New("unable to get tags for packages by package ids")
						for i := 0; i < len(ids); i++ {
							errs[i] = e
						}
						return
					}

					var idMap = make(map[graphtypes.ID][]graphtypes.ID, len(dbTags))
					for _, tag := range dbTags {
						id := graphtypes.ID(tag.PackageID)
						idMap[id] = append(idMap[id], graphtypes.ID(tag.TagID))
					}
					for i, id := range ids {
						packageTags[i] = idMap[id]
					}

					return
				},
			},
			TagById: TagLoader{
				maxBatch: 100,
				wait:     time.Millisecond,
				fetch: func(ids []graphtypes.ID) (tags []*model.Tag, errs []error) {
					errs = make([]error, len(ids))
					tags = make([]*model.Tag, len(ids))

					var dbTags []*database.Tag
					if err := util.DBFromContext(ctx).Find(&dbTags, ids).Error; err != nil {
						log.Printf("Error: %v", err)
						e := errors.New("unable to get tags by tag ids")
						for i := 0; i < len(ids); i++ {
							errs[i] = e
						}
						return
					}

					var idMap = make(map[graphtypes.ID]*model.Tag, len(dbTags))
					conv := generated.ConverterDBImpl{}
					for _, tag := range dbTags {
						mTag := conv.ConvertTag(*tag)
						idMap[graphtypes.ID(tag.ID)] = &mTag
					}
					for i, id := range ids {
						tags[i] = idMap[id]
					}
					return
				},
			},
			ReleasesByPackage: ReleasesByPackageLoader{
				maxBatch: 100,
				wait:     time.Millisecond,
				fetch: func(ids []graphtypes.ID) (releases [][]*model.Release, errs []error) {
					releases = make([][]*model.Release, len(ids))
					errs = make([]error, len(ids))

					var dbReleases []*database.Release
					if err := util.DBFromContext(ctx).Where("package_id IN ?", ids).Find(&dbReleases).Error; err != nil {
						log.Printf("Error: %v", err)
						e := errors.New("unable to get releases by package ids")
						for i := 0; i < len(ids); i++ {
							errs[i] = e
						}
						return
					}

					var idMap = make(map[graphtypes.ID][]*model.Release, len(dbReleases))
					conv := generated.ConverterDBImpl{}
					for _, release := range dbReleases {
						id := graphtypes.ID(release.PackageID)
						r := conv.ConvertRelease(*release)
						idMap[id] = append(idMap[id], &r)
					}
					for i, id := range ids {
						releases[i] = idMap[id]
					}
					return
				},
			},
			PackageById: PackageByIdLoader{
				maxBatch: 100,
				wait:     time.Millisecond,
				fetch: func(ids []graphtypes.ID) (packages []*model.Package, errs []error) {
					packages = make([]*model.Package, len(ids))
					errs = make([]error, len(ids))

					var dbPackages []*database.Package
					if err := util.DBFromContext(ctx).Find(&dbPackages, ids).Error; err != nil {
						log.Printf("Error: %v", err)
						e := errors.New("unable to get packages by package ids")
						for i := 0; i < len(ids); i++ {
							errs[i] = e
						}
						return
					}

					var idMap = make(map[graphtypes.ID]*model.Package, len(dbPackages))
					conv := generated.ConverterDBImpl{}
					for _, pack := range dbPackages {
						idMap[graphtypes.ID(pack.ID)] = conv.ConvertPackageP(pack)
					}
					for i, id := range ids {
						packages[i] = idMap[id]
					}
					return
				},
			},
			PackagesByTag: PackagesByTagLoader{
				maxBatch: 100,
				wait:     time.Millisecond,
				fetch: func(ids []graphtypes.ID) (tagPackages [][]graphtypes.ID, errs []error) {
					errs = make([]error, len(ids))
					tagPackages = make([][]graphtypes.ID, len(ids))

					var dbTags []*database.PackageTag
					if err := util.DBFromContext(ctx).Where("tag_id IN ?", ids).Find(&dbTags).Error; err != nil {
						log.Printf("Error: %v", err)
						e := errors.New("unable to get packages for tags by tag ids")
						for i := 0; i < len(ids); i++ {
							errs[i] = e
						}
						return
					}

					var idMap = make(map[graphtypes.ID][]graphtypes.ID, len(dbTags))
					for _, tag := range dbTags {
						id := graphtypes.ID(tag.TagID)
						idMap[id] = append(idMap[id], graphtypes.ID(tag.PackageID))
					}
					for i, id := range ids {
						tagPackages[i] = idMap[id]
					}

					return
				},
			},
			PackagesByUser: PackagesByUserLoader{
				maxBatch: 100,
				wait:     time.Millisecond,
				fetch: func(ids []graphtypes.ID) (packages [][]*model.Package, errs []error) {
					packages = make([][]*model.Package, len(ids))
					errs = make([]error, len(ids))

					var dbPackages []*database.Package
					if err := util.DBFromContext(ctx).Where("package_creator_id IN ?", ids).Find(&dbPackages).Error; err != nil {
						log.Printf("Error: %v", err)
						e := errors.New("unable to get packages by user ids")
						for i := 0; i < len(ids); i++ {
							errs[i] = e
						}
						return
					}

					var idMap = make(map[graphtypes.ID][]*model.Package, len(dbPackages))
					conv := generated.ConverterDBImpl{}
					for _, pack := range dbPackages {
						id := graphtypes.ID(pack.CreatorID)
						idMap[id] = append(idMap[id], conv.ConvertPackageP(pack))
					}
					for i, id := range ids {
						packages[i] = idMap[id]
					}
					return
				},
			},
		})
		c.SetRequest(c.Request().WithContext(nCtx))
		return handlerFunc(c)
	}
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey{}).(*Loaders)
}
