package auth

import (
	"FINRepository/Database"
	"context"
)

func AuthorizeViewPackage(ctx context.Context, p *Database.Package) bool {
	if p.Verified {
		return true
	}

	user := UserFromCtx(ctx)
	if user == nil {
		return false
	}

	authorizer := AuthorizerFromCtx(ctx)
	success, _ := authorizer.Authorize(ctx, p, user, "view_all")
	return success
}
