package perm

import "context"

var Global AuthorizableGeneric = AuthorizableGeneric{"global", "0"}

type Authorizable interface {
	GetType() string
	GetID() string
}

type AuthorizableGeneric struct {
	Type string
	ID   string
}

func (a *AuthorizableGeneric) GetType() string {
	return a.Type
}

func (a *AuthorizableGeneric) GetID() string {
	return a.ID
}

type Authorizer interface {
	Authorize(ctx context.Context, resource Authorizable, subject Authorizable, permission string) (bool, error)
	AddRelation(ctx context.Context, resource Authorizable, subject Authorizable, relation string) (string, error)
	RemoveRelation(ctx context.Context, resource Authorizable, subject Authorizable, relation string) (string, error)
}

func CtxWithAuthorizer(ctx context.Context, authorizer Authorizer) context.Context {
	return context.WithValue(ctx, "authorizer", authorizer)
}

func AuthorizerFromCtx(ctx context.Context) Authorizer {
	return ctx.Value("authorizer").(Authorizer)
}
