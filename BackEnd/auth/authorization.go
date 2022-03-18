package auth

import "context"

type Authorizable interface {
	GetType() string
	GetID() string
}

type Authorizer interface {
	Authorize(ctx context.Context, resource Authorizable, subject Authorizable, permission string) (bool, error)
	Permit(ctx context.Context, resource Authorizable, subject Authorizable, relation string) error
	RemovePermit(ctx context.Context, resource Authorizable, subject Authorizable, relation string) error
}

func CtxWithAuthorizer(ctx context.Context, authorizer Authorizer) context.Context {
	return context.WithValue(ctx, "authorizer", authorizer)
}

func AuthorizerFromCtx(ctx context.Context) Authorizer {
	return ctx.Value("authorizer").(Authorizer)
}
