package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	generated1 "FINRepository/Convert/generated"
	"FINRepository/Database"
	"FINRepository/Util"
	"FINRepository/dataloader"
	"FINRepository/graph/generated"
	"FINRepository/graph/model"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/99designs/gqlgen/graphql"
)

func (r *mutationResolver) CreatePackage(ctx context.Context, packageArg model.NewPackage) (*model.Package, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdatePackage(ctx context.Context, packageArg model.UpdatePackage) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeletePackage(ctx context.Context, packageID uint64) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) NewRelease(ctx context.Context, release model.NewRelease) (*model.Release, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateRelease(ctx context.Context, release model.UpdateRelease) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteRelease(ctx context.Context, releaseID uint64) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteUser(ctx context.Context, userID uint64) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateTag(ctx context.Context, tag model.NewTag) (*model.Tag, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateTag(ctx context.Context, tag model.UpdateTag) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteTag(ctx context.Context, tagID uint64) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *packageResolver) Creator(ctx context.Context, obj *model.Package) (*model.User, error) {
	return dataloader.For(ctx).UserById.Load(obj.Creator.ID)
}

func (r *packageResolver) Tags(ctx context.Context, obj *model.Package) ([]*model.Tag, error) {
	tagIds, err := dataloader.For(ctx).TagsByPackage.Load(obj.ID)
	if err != nil {
		return nil, err
	}
	tags, _ := dataloader.For(ctx).TagById.LoadAll(tagIds)
	return tags, nil
}

func (r *packageResolver) Releases(ctx context.Context, obj *model.Package) ([]*model.Release, error) {
	return dataloader.For(ctx).ReleasesByPackage.Load(obj.ID)
}

func (r *queryResolver) ListPackages(ctx context.Context, page int, count int) ([]*model.Package, error) {
	var fieldMap = map[string]string{
		"id":          "package_id",
		"name":        "package_name",
		"displayName": "package_displayname",
		"description": "package_description",
		"sourceLink":  "package_sourcelink",
		"verified":    "package_verified",
		"creator":     "package_creator_id",
		"tags":        "package_id",
		"releases":    "package_id",
	}
	colFields := graphql.CollectFieldsCtx(ctx, nil)

	var query = Util.DBFromContext(ctx).Scopes(Util.Paginate(page, count))

	fields := make([]string, len(colFields))
	for i, field := range colFields {
		fields[i] = fieldMap[field.Name]
	}
	query = query.Select(fields)

	var packages []*Database.Package
	if err := query.Find(&packages).Error; err != nil {
		return nil, errors.New("unable to get packages")
	}

	var packs = make([]*model.Package, len(packages))
	var conv = generated1.ConverterImpl{}
	for i, pack := range packages {
		p := conv.ConvertPackage(*pack)
		packs[i] = &p
	}

	return packs, nil
}

func (r *queryResolver) GetPackagesByID(ctx context.Context, ids []uint64) ([]*model.Package, error) {
	var fieldMap = map[string]string{
		"id":          "package_id",
		"name":        "package_name",
		"displayName": "package_displayname",
		"description": "package_description",
		"sourceLink":  "package_sourcelink",
		"verified":    "package_verified",
		"creator":     "package_creator_id",
		"tags":        "package_id",
		"releases":    "package_id",
	}
	colFields := graphql.CollectFieldsCtx(ctx, nil)

	var query = Util.DBFromContext(ctx)

	fields := make([]string, len(colFields))
	for i, field := range colFields {
		fields[i] = fieldMap[field.Name]
	}
	query = query.Select(fields)

	var packages []*Database.Package
	if err := query.Find(&packages, ids).Error; err != nil {
		return nil, errors.New("unable to get packages")
	}

	var idMap = make(map[uint64]*model.Package, len(packages))
	var conv = generated1.ConverterImpl{}
	for _, pack := range packages {
		p := conv.ConvertPackage(*pack)
		idMap[pack.ID] = &p
	}
	var packs = make([]*model.Package, len(ids))
	for i, id := range ids {
		packs[i] = idMap[id]
	}

	return packs, nil
}

func (r *queryResolver) GetUsersByID(ctx context.Context, ids []uint64) ([]*model.User, error) {
	var fieldMap = map[string]string{
		"id":       "user_id",
		"name":     "user_name",
		"bio":      "user_bio",
		"admin":    "user_admin",
		"email":    "user_email",
		"verified": "user_verified",
		"packages": "user_id",
	}
	colFields := graphql.CollectFieldsCtx(ctx, nil)

	var query = Util.DBFromContext(ctx)

	fields := make([]string, len(colFields))
	for i, field := range colFields {
		fields[i] = fieldMap[field.Name]
	}
	query = query.Select(fields)

	var dbUsers []*Database.User
	if err := query.Find(&dbUsers, ids).Error; err != nil {
		return nil, errors.New("unable to get idMap")
	}

	var idMap = make(map[uint64]*model.User, len(dbUsers))
	var conv = generated1.ConverterImpl{}
	for _, user := range dbUsers {
		u := conv.ConvertUser(*user)
		idMap[user.ID] = &u
	}
	var packs = make([]*model.User, len(ids))
	for i, id := range ids {
		packs[i] = idMap[id]
	}
	return packs, nil
}

func (r *queryResolver) GetAllTags(ctx context.Context) ([]*model.Tag, error) {
	dbTags, err := Database.GetTags(Util.DBFromContext(ctx))
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, errors.New("unable to get all tags")
	}
	conv := generated1.ConverterImpl{}
	return conv.ConvertTagPA(*dbTags), nil
}

func (r *releaseResolver) Package(ctx context.Context, obj *model.Release) (*model.Package, error) {
	return dataloader.For(ctx).PackageById.Load(obj.ID)
}

func (r *tagResolver) Packages(ctx context.Context, obj *model.Tag) ([]*model.Package, error) {
	packagesId, err := dataloader.For(ctx).PackagesByTag.Load(obj.ID)
	if err != nil {
		return nil, err
	}
	packages, _ := dataloader.For(ctx).PackageById.LoadAll(packagesId)
	return packages, nil
}

func (r *userResolver) Packages(ctx context.Context, obj *model.User) ([]*model.Package, error) {
	return dataloader.For(ctx).PackagesByUser.Load(obj.ID)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Package returns generated.PackageResolver implementation.
func (r *Resolver) Package() generated.PackageResolver { return &packageResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Release returns generated.ReleaseResolver implementation.
func (r *Resolver) Release() generated.ReleaseResolver { return &releaseResolver{r} }

// Tag returns generated.TagResolver implementation.
func (r *Resolver) Tag() generated.TagResolver { return &tagResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type packageResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type releaseResolver struct{ *Resolver }
type tagResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
