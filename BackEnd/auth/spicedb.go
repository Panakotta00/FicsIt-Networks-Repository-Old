package auth

import (
	"context"
	_ "embed"
	pb "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/authzed/grpcutil"
	"log"
)

//go:embed schema.sdb
var schema string

type Authorizer_SpiceDB struct {
	client *authzed.Client
}

func (*Authorizer_SpiceDB) Authorize(Resource Authorizeable, Subject Authorizeable, Permission string) {
}
func (*Authorizer_SpiceDB) Permit(Resource Authorizeable, Subject Authorizeable, Relation string) {}
func (*Authorizer_SpiceDB) RemovePermit(Resource Authorizeable, Subject Authorizeable, Relation string) {
}

func NewSpiceDBAuthorizer(host string, token string) (*Authorizer_SpiceDB, error) {
	auth := new(Authorizer_SpiceDB)

	client, err := authzed.NewClient(
		host,
		// grpcutil.WithBearerToken("t_your_token_here_1234567deadbeef"),
		// grpcutil.WithSystemCerts(grpcutil.VerifyCA),
		grpcutil.WithInsecureBearerToken(token),
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
