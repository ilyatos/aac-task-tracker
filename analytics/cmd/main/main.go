package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/ilyatos/aac-task-tracker/analytics/internal/consumer/transactions_lifecycle"

	add_task_prices_command "github.com/ilyatos/aac-task-tracker/analytics/internal/app/command/add_task_prices"
	complete_task_command "github.com/ilyatos/aac-task-tracker/analytics/internal/app/command/complete_task"
	create_task_command "github.com/ilyatos/aac-task-tracker/analytics/internal/app/command/create_task"
	create_transaction_command "github.com/ilyatos/aac-task-tracker/analytics/internal/app/command/create_transaction"
	"github.com/ilyatos/aac-task-tracker/analytics/internal/consumer/prices_stream"
	"github.com/ilyatos/aac-task-tracker/analytics/internal/consumer/tasks_lifecycle"

	"github.com/ilyatos/aac-task-tracker/analytics/internal/consumer/tasks_stream"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	"github.com/ilyatos/aac-task-tracker/analytics/gen/analytics"
	create_user_command "github.com/ilyatos/aac-task-tracker/analytics/internal/app/command/create_user"
	"github.com/ilyatos/aac-task-tracker/analytics/internal/app/helper/jwt"
	"github.com/ilyatos/aac-task-tracker/analytics/internal/app/repository"
	"github.com/ilyatos/aac-task-tracker/analytics/internal/consumer/users_stream"
	"github.com/ilyatos/aac-task-tracker/analytics/internal/context_key"
	"github.com/ilyatos/aac-task-tracker/analytics/internal/database"
	"github.com/ilyatos/aac-task-tracker/analytics/internal/middleware"
	"github.com/ilyatos/aac-task-tracker/analytics/internal/rpc"
	"github.com/ilyatos/aac-task-tracker/analytics/internal/rpc/get_most_expensive_task"
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
	transactionRepository := repository.NewTransactionRepository(db)

	// Commands and queries
	createUserCommand := create_user_command.New(userRepository)
	createTaskCommand := create_task_command.New(taskRepository)
	addTaskPricesCommand := add_task_prices_command.New(taskRepository)
	completeTaskCommand := complete_task_command.New(taskRepository)
	createTransactionCommand := create_transaction_command.New(transactionRepository)

	// Register handlers
	handler := rpc.New(
		get_most_expensive_task.New().Handle,
	)

	// Broker consumers
	usersStreamConsumer := users_stream.New(createUserCommand)
	go usersStreamConsumer.Consume()

	tasksStreamConsumer := tasks_stream.New(createTaskCommand)
	go tasksStreamConsumer.Consume()

	pricesStreamConsumer := prices_stream.New(addTaskPricesCommand)
	go pricesStreamConsumer.Consume()

	tasksConsumer := tasks_lifecycle.New(completeTaskCommand)
	go tasksConsumer.Consume()

	transactionConsumer := transactions_lifecycle.New(createTransactionCommand)
	go transactionConsumer.Consume()

	// gRPC server
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middleware.Log(), middleware.JWT(jwtManager)),
	)
	analytics.RegisterAnalyticsServer(srv, handler)
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
	err = analytics.RegisterAnalyticsHandler(context.Background(), mux, conn)
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
