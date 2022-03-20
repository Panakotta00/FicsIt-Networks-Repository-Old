package auth

import (
	"FINRepository/auth/perm"
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

func AuthorizableToObjRef(auth perm.Authorizable) *pb.ObjectReference {
	return &pb.ObjectReference{
		ObjectType: auth.GetType(),
		ObjectId:   auth.GetID(),
	}
}

func AuthorizableToSubRef(auth perm.Authorizable) *pb.SubjectReference {
	return &pb.SubjectReference{
		Object: &pb.ObjectReference{
			ObjectType: auth.GetType(),
			ObjectId:   auth.GetID(),
		},
	}
}

func (db *Authorizer_SpiceDB) Authorize(ctx context.Context, resource perm.Authorizable, subject perm.Authorizable, permission string) (bool, error) {
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

func (db *Authorizer_SpiceDB) AddRelation(ctx context.Context, resource perm.Authorizable, subject perm.Authorizable, relation string) (string, error) {
	resp, err := db.client.WriteRelationships(ctx, &pb.WriteRelationshipsRequest{
		Updates: []*pb.RelationshipUpdate{
			{
				Operation: pb.RelationshipUpdate_OPERATION_TOUCH,
				Relationship: &pb.Relationship{
					Resource: AuthorizableToObjRef(resource),
					Relation: relation,
					Subject:  AuthorizableToSubRef(subject),
				},
			},
		},
	})
	if err != nil {
		return "", err
	}
	return resp.GetWrittenAt().GetToken(), nil
}

func (db *Authorizer_SpiceDB) RemoveRelation(ctx context.Context, resource perm.Authorizable, subject perm.Authorizable, relation string) (string, error) {
	resp, err := db.client.WriteRelationships(ctx, &pb.WriteRelationshipsRequest{
		Updates: []*pb.RelationshipUpdate{
			{
				Operation: pb.RelationshipUpdate_OPERATION_DELETE,
				Relationship: &pb.Relationship{
					Resource: AuthorizableToObjRef(resource),
					Relation: relation,
					Subject:  AuthorizableToSubRef(subject),
				},
			},
		},
	})
	if err != nil {
		return "", err
	}
	return resp.GetWrittenAt().GetToken(), nil
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
