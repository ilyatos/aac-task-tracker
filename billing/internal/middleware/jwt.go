package middleware

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/ilyatos/aac-task-tracker/billing/internal/app/helper/jwt"
	"github.com/ilyatos/aac-task-tracker/billing/internal/context_key"
)

func JWT(manager *jwt.Manager) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		log.Printf("auth attempt to %s: '%s'\n", info.FullMethod, req)

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		authValues := md.Get(context_key.AuthKey)
		if len(authValues) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		accessToken := authValues[0]

		userClaims, err := manager.ParseAccessToken(accessToken)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, fmt.Sprintf("token is malformed: %s", err))
		}

		ctx = context.WithValue(ctx, context_key.AuthKey, userClaims)

		resp, err = handler(ctx, req)

		log.Printf("success auth attempt to %s: '%s'\n", info.FullMethod, req)

		return
	}
}
