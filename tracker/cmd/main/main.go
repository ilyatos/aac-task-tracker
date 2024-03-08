package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/ilyatos/aac-task-tracker/tracker/internal/context_key"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	"github.com/ilyatos/aac-task-tracker/tracker/gen/tracker"
	complete_task_command "github.com/ilyatos/aac-task-tracker/tracker/internal/app/command/complete_task"
	create_task_command "github.com/ilyatos/aac-task-tracker/tracker/internal/app/command/create_task"
	create_user_command "github.com/ilyatos/aac-task-tracker/tracker/internal/app/command/create_user"
	reassign_tasks_command "github.com/ilyatos/aac-task-tracker/tracker/internal/app/command/reassign_tasks"
	update_user_command "github.com/ilyatos/aac-task-tracker/tracker/internal/app/command/update_user"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/app/helper/jwt"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/app/repository"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/consumer/users_stream"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/database"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/middleware"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/rpc"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/rpc/complete_task"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/rpc/create_task"
	"github.com/ilyatos/aac-task-tracker/tracker/internal/rpc/reassign_tasks"
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

	// Repositories
	userRepository := repository.NewUserRepository(db)
	taskRepository := repository.NewTaskRepository(db)

	// Commands and queries
	createTaskCommand := create_task_command.New(userRepository, taskRepository)
	reassignTasksCommand := reassign_tasks_command.New(userRepository, taskRepository)
	markTaskAsDoneCommand := complete_task_command.New(taskRepository)
	createUserCommand := create_user_command.New(userRepository)
	updateUserCommand := update_user_command.New(userRepository)

	// Register handlers
	handler := rpc.New(
		create_task.New(createTaskCommand).Handle,
		reassign_tasks.New(reassignTasksCommand).Handle,
		complete_task.New(markTaskAsDoneCommand).Handle,
	)

	// Broker consumers
	usersStreamConsumer := users_stream.NewConsumer(createUserCommand, updateUserCommand)
	go usersStreamConsumer.Consume()

	// gRPC server
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middleware.Log(), middleware.JWT(jwtManager)),
	)
	tracker.RegisterTrackerServer(srv, handler)
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
			return metadata.Pairs(context_key.AuthKey, header)
		}),
	)
	err = tracker.RegisterTrackerHandler(context.Background(), mux, conn)
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
