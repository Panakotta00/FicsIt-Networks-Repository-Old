package auth

import "context"

type Authorizeable interface {
}

type Authorizer interface {
	Authorize(Resource Authorizeable, Subject Authorizeable, Permission string)
	Permit(Resource Authorizeable, Subject Authorizeable, Permission string)
	RemovePermit(Resource Authorizeable, Subject Authorizeable, Permission string)
}

func CtxWithAuthorizer(ctx context.Context, authorizer Authorizer) context.Context {
	return context.WithValue(ctx, "authorizer", authorizer)
}

func AuthorizerFromCtx(ctx context.Context) Authorizer {
	return ctx.Value("authorizer").(Authorizer)
}
