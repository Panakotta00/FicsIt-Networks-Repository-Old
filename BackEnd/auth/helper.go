package auth

import (
	"FINRepository/auth/perm"
	"FINRepository/database"
	"FINRepository/graph/model"
	"context"
)

func AuthorizeViewAll(ctx context.Context, obj perm.Authorizable) bool {
	user := database.UserFromCtx(ctx)
	if user == nil {
		return false
	}

	authorizer := perm.AuthorizerFromCtx(ctx)
	success, _ := authorizer.Authorize(ctx, obj, user, "view_all")
	return success
}

type Verifiable interface {
	perm.Authorizable
	IsVerified() bool
}

func AuthorizeVerification(ctx context.Context, obj Verifiable) bool {
	if obj.IsVerified() {
		return true
	}

	return AuthorizeViewAll(ctx, obj)
}

func AuthorizeVerifications[V Verifiable](ctx context.Context, objs []V) []V {
	filtered := make([]V, len(objs))
	i := 0
	for _, obj := range objs {
		if AuthorizeVerification(ctx, obj) {
			filtered[i] = obj
			i = i + 1
		}
	}
	return filtered[:i]
}

func AuthorizeManageTags(ctx context.Context) bool {
	user := database.UserFromCtx(ctx)
	if user == nil {
		return false
	}

	authorizer := perm.AuthorizerFromCtx(ctx)
	success, _ := authorizer.Authorize(ctx, &perm.Global, user, "edit_tags")
	return success
}

func AuthorizeTagModel(ctx context.Context, r *model.Tag) bool {
	if r.Verified {
		return true
	}

	return database.AuthorizeCtx(ctx, &perm.Global, "edit_tags")
}

func AuthorizeTagModels(ctx context.Context, tags []*model.Tag) []*model.Tag {
	filteredTags := make([]*model.Tag, len(tags))
	i := 0
	for _, tag := range tags {
		if AuthorizeTagModel(ctx, tag) {
			filteredTags[i] = tag
			i = i + 1
		}
	}
	return filteredTags[:i]
}
