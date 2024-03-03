package middleware

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/ilyatos/aac-task-tracker/auth/internal/app/helper/jwt"
)

var endpointsWithoutAuth = map[string]struct{}{
	"/auth.Auth/LogIn":  {},
	"/auth.Auth/SignUp": {},
}

const AuthKeyName = "auth"

func JWT(manager *jwt.Manager) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		log.Printf("auth attempt to %s: '%s'\n", info.FullMethod, req)

		_, ok := endpointsWithoutAuth[info.FullMethod]
		if ok {
			log.Printf("auth not required to %s: '%s'\n", info.FullMethod, req)
			resp, err = handler(ctx, req)
			return
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
		}

		authValues := md.Get(AuthKeyName)
		if len(authValues) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
		}

		accessToken := authValues[0]

		_, err = manager.ParseAccessToken(accessToken)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, fmt.Sprintf("token is malformed: %s", err))
		}

		resp, err = handler(ctx, req)

		log.Printf("success auth attempt to %s: '%s'\n", info.FullMethod, req)

		return
	}
}
