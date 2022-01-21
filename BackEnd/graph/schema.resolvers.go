package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"FINRepository/Convert/generated"
	"FINRepository/Database"
	"FINRepository/Util"
	"FINRepository/dataloader"
	generated1 "FINRepository/graph/generated"
	"FINRepository/graph/model"
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
)

func (r *packageResolver) Creator(ctx context.Context, obj *model.Package) (*model.User, error) {
	return dataloader.For(ctx).UserById.Load(obj.Creator.ID)
}

func (r *queryResolver) Packages(ctx context.Context) ([]*model.Package, error) {
	var fieldMap map[string]string = map[string]string{
		"id":          "package_id",
		"name":        "package_name",
		"displayName": "package_displayname",
		"description": "package_description",
		"sourceLink":  "package_sourcelink",
		"verified":    "package_verified",
		"creator":     "package_creator_id",
	}
	colFields := graphql.CollectFieldsCtx(ctx, nil)
	var fields []string = make([]string, len(colFields))
	for i, field := range colFields {
		fields[i] = fieldMap[field.Name]
	}

	var packages []*Database.Package
	if err := Util.DBFromContext(ctx).Scopes(Util.Paginate(0, 50)).Select(fields).Find(&packages).Error; err != nil {
		return nil, errors.New("unable to get packages")
	}

	var packs = make([]*model.Package, len(packages))
	var conv = generated.ConverterImpl{}
	for i, pack := range packages {
		p := conv.ConvertPackage(*pack)
		packs[i] = &p
	}

	return packs, nil
}

// Package returns generated1.PackageResolver implementation.
func (r *Resolver) Package() generated1.PackageResolver { return &packageResolver{r} }

// Query returns generated1.QueryResolver implementation.
func (r *Resolver) Query() generated1.QueryResolver { return &queryResolver{r} }

type packageResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
