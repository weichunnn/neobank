package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/weichunnn/neobank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationType   = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}
	authHeader := values[0]
	fields := strings.Fields(authHeader)

	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationType {
		return nil, fmt.Errorf("unsupported authorization type")
	}

	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token: %s", err)
	}

	return payload, nil
}
