package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	generated1 "FINRepository/Convert/generated"
	"FINRepository/Database"
	"FINRepository/Util"
	"FINRepository/dataloader"
	"FINRepository/graph/generated"
	"FINRepository/graph/graphtypes"
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

func (r *mutationResolver) DeletePackage(ctx context.Context, packageID graphtypes.ID) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) NewRelease(ctx context.Context, release model.NewRelease) (*model.Release, error) {
	rel, err := Database.CreateRelease(ctx, Database.ID(release.PackageID), release.Name, release.Description, release.SourceLink, release.Version, release.FinVersion)
	if err != nil {
		return nil, errors.New("Unable to create new release")
	}
	conv := generated1.ConverterDBImpl{}
	return conv.ConvertReleaseP(rel), nil
}

func (r *mutationResolver) UpdateRelease(ctx context.Context, release model.UpdateRelease) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteRelease(ctx context.Context, releaseID graphtypes.ID) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateTag(ctx context.Context, tag model.NewTag) (*model.Tag, error) {
	dbTag, err := Database.CreateTag(ctx, tag.Name, tag.Description)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, errors.New("unable to create tag")
	}
	var conv = &generated1.ConverterDBImpl{}
	return conv.ConvertTagP(dbTag), nil
}

func (r *mutationResolver) UpdateTag(ctx context.Context, tag model.UpdateTag) (bool, error) {
	return Database.UpdateTag(ctx, Database.ID(tag.ID), tag.Name, tag.Description), nil
}

func (r *mutationResolver) DeleteTag(ctx context.Context, tagID graphtypes.ID) (bool, error) {
	return Database.DeleteTag(ctx, Database.ID(tagID)), nil
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
		"verified":    "package_verified, package_creator_id",
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
	var conv = generated1.ConverterDBImpl{}
	for i, pack := range packages {
		p := conv.ConvertPackage(*pack)
		packs[i] = &p
	}

	return packs, nil
}

func (r *queryResolver) GetPackagesByID(ctx context.Context, ids []graphtypes.ID) ([]*model.Package, error) {
	var fieldMap = map[string]string{
		"id":          "package_id",
		"name":        "package_name",
		"displayName": "package_displayname",
		"description": "package_description",
		"sourceLink":  "package_sourcelink",
		"verified":    "package_verified, package_creator_id",
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

	var idMap = make(map[graphtypes.ID]*model.Package, len(packages))
	var conv = generated1.ConverterDBImpl{}
	for _, pack := range packages {
		p := conv.ConvertPackage(*pack)
		idMap[graphtypes.ID(pack.ID)] = &p
	}
	var packs = make([]*model.Package, len(ids))
	for i, id := range ids {
		packs[i] = idMap[id]
	}

	return packs, nil
}

func (r *queryResolver) GetUsersByID(ctx context.Context, ids []graphtypes.ID) ([]*model.User, error) {
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

	var idMap = make(map[graphtypes.ID]*model.User, len(dbUsers))
	var conv = generated1.ConverterDBImpl{}
	for _, user := range dbUsers {
		u := conv.ConvertUser(*user)
		idMap[graphtypes.ID(user.ID)] = &u
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
	conv := generated1.ConverterDBImpl{}
	tags := conv.ConvertTagPA(*dbTags)
	return tags, nil
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
