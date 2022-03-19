package auth

import (
	"context"
	_ "embed"
	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"google.golang.org/grpc"
	"log"
)

//go:embed schema.sdb
var schema string

type Authorizer_SpiceDB struct {
	client *authzed.Client
}

func AuthorizableToObjRef(auth Authorizable) *pb.ObjectReference {
	return &pb.ObjectReference{
		ObjectType: auth.GetType(),
		ObjectId:   auth.GetID(),
	}
}

func AuthorizableToSubRef(auth Authorizable) *pb.SubjectReference {
	return &pb.SubjectReference{
		Object: &pb.ObjectReference{
			ObjectType: auth.GetType(),
			ObjectId:   auth.GetID(),
		},
	}
}

func (db *Authorizer_SpiceDB) Authorize(ctx context.Context, resource Authorizable, subject Authorizable, permission string) (bool, error) {
	resp, err := db.client.CheckPermission(ctx, &pb.CheckPermissionRequest{
		Resource:   AuthorizableToObjRef(resource),
		Permission: permission,
		Subject:    AuthorizableToSubRef(subject),
	})
	if err != nil {
		log.Printf("Authorization error: %v", err)
		return false, nil
	}
	return resp.Permissionship == pb.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION, nil
}

func (*Authorizer_SpiceDB) AddRelation(ctx context.Context, resource Authorizable, subject Authorizable, relation string) error {
	return nil
}

func (*Authorizer_SpiceDB) RemoveRelation(ctx context.Context, resource Authorizable, subject Authorizable, relation string) error {
	return nil
}

func NewSpiceDBAuthorizer(host string, token string) (*Authorizer_SpiceDB, error) {
	auth := new(Authorizer_SpiceDB)

	client, err := authzed.NewClient(
		host,
		// grpcutil.WithBearerToken("t_your_token_here_1234567deadbeef"),
		// grpcutil.WithSystemCerts(grpcutil.VerifyCA),
		grpcutil.WithInsecureBearerToken(token),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("unable to initialize client: %s", err)
	}

	request := &pb.WriteSchemaRequest{Schema: schema}
	_, err = client.WriteSchema(context.Background(), request)
	if err != nil {
		return nil, err
	}

	auth.client = client

	return auth, nil
}
