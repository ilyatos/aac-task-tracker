package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	"github.com/ilyatos/aac-task-tracker/auth/gen/auth"
	change_role_rommand "github.com/ilyatos/aac-task-tracker/auth/internal/app/command/change_role"
	login_command "github.com/ilyatos/aac-task-tracker/auth/internal/app/command/login"
	singup_command "github.com/ilyatos/aac-task-tracker/auth/internal/app/command/singup"
	"github.com/ilyatos/aac-task-tracker/auth/internal/app/helper/jwt"
	"github.com/ilyatos/aac-task-tracker/auth/internal/app/helper/password"
	"github.com/ilyatos/aac-task-tracker/auth/internal/app/repository"
	"github.com/ilyatos/aac-task-tracker/auth/internal/database"
	"github.com/ilyatos/aac-task-tracker/auth/internal/middleware"
	"github.com/ilyatos/aac-task-tracker/auth/internal/rpc"
	"github.com/ilyatos/aac-task-tracker/auth/internal/rpc/change_role"
	"github.com/ilyatos/aac-task-tracker/auth/internal/rpc/login"
	"github.com/ilyatos/aac-task-tracker/auth/internal/rpc/signup"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = godotenv.Load()
	if err != nil {
		return
	}

	// Database connection
	db, err := database.Connect()
	if err != nil {
		return
	}

	// Helpers
	jwtManager := jwt.New()
	passwordHasher := password.New()

	// Repositories
	userRepository := repository.NewUserRepository(db)

	// Commands and queries
	signupCommand := singup_command.New(userRepository, passwordHasher)
	loginCommand := login_command.New(userRepository, passwordHasher, jwtManager)
	changeRoleCommand := change_role_rommand.New(userRepository)

	// Register handlers
	handler := rpc.New(
		signup.New(signupCommand).Handle,
		login.New(loginCommand).Handle,
		change_role.New(changeRoleCommand).Handle,
	)

	// gRPC server
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middleware.Log(), middleware.JWT(jwtManager)),
	)
	auth.RegisterAuthServer(srv, handler)
	reflection.Register(srv)

	l, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		return
	}

	log.Println("Serving gRPC server")

	go func() {
		log.Fatalln(srv.Serve(l))
	}()

	// HTTP proxy to the gRPC server
	conn, err := grpc.Dial("0.0.0.0:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	defer conn.Close()

	mux := runtime.NewServeMux(
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			header := request.Header.Get("Authorization")
			md := metadata.Pairs(middleware.AuthKeyName, header)
			return md
		}),
	)
	err = auth.RegisterAuthHandler(context.Background(), mux, conn)
	if err != nil {
		return
	}

	gwServer := &http.Server{
		Addr:    "0.0.0.0:9001",
		Handler: mux,
	}

	log.Println("Serving gRPC HTTP proxy")
	log.Fatalln(gwServer.ListenAndServe())
}
