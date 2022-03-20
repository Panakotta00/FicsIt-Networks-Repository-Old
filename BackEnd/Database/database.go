package Database

import (
	"FINRepository/auth/perm"
	"context"
)

type ID int64

func UserFromCtx(ctx context.Context) *User {
	val := ctx.Value("auth")
	if val != nil {
		return val.(*User)
	}
	return nil
}

func AuthorizeCtx(ctx context.Context, resource perm.Authorizable, permission string) bool {
	user := UserFromCtx(ctx)
	if user == nil {
		return false
	}

	authorizer := perm.AuthorizerFromCtx(ctx)
	success, _ := authorizer.Authorize(ctx, resource, user, permission)
	return success
}
