package middleware

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func Log() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		log.Printf("req to %s: '%s'\n", info.FullMethod, req)
		resp, err = handler(ctx, req)
		log.Printf("resp withto %s: %s; err: %s\n", info.FullMethod, resp, err)

		return
	}
}
