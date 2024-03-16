package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-co-op/gocron/v2"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	"github.com/ilyatos/aac-task-tracker/billing/gen/billing"
	create_credit_transaction_command "github.com/ilyatos/aac-task-tracker/billing/internal/app/command/create_credit_transaction"
	create_debit_transaction_command "github.com/ilyatos/aac-task-tracker/billing/internal/app/command/create_debit_transaction"
	create_payout_transaction_command "github.com/ilyatos/aac-task-tracker/billing/internal/app/command/create_payout_transaction"
	create_task_command "github.com/ilyatos/aac-task-tracker/billing/internal/app/command/create_task"
	create_transaction_command "github.com/ilyatos/aac-task-tracker/billing/internal/app/command/create_transaction"
	create_user_command "github.com/ilyatos/aac-task-tracker/billing/internal/app/command/create_user"
	update_task_command "github.com/ilyatos/aac-task-tracker/billing/internal/app/command/update_task"
	update_user_command "github.com/ilyatos/aac-task-tracker/billing/internal/app/command/update_user"
	"github.com/ilyatos/aac-task-tracker/billing/internal/app/helper/jwt"
	get_user_balance_query "github.com/ilyatos/aac-task-tracker/billing/internal/app/query/get_user_balance"
	"github.com/ilyatos/aac-task-tracker/billing/internal/app/repository"
	"github.com/ilyatos/aac-task-tracker/billing/internal/consumer/tasks_lifecycle"
	"github.com/ilyatos/aac-task-tracker/billing/internal/consumer/tasks_stream"
	"github.com/ilyatos/aac-task-tracker/billing/internal/consumer/users_stream"
	"github.com/ilyatos/aac-task-tracker/billing/internal/context_key"
	"github.com/ilyatos/aac-task-tracker/billing/internal/database"
	"github.com/ilyatos/aac-task-tracker/billing/internal/middleware"
	"github.com/ilyatos/aac-task-tracker/billing/internal/rpc"
	"github.com/ilyatos/aac-task-tracker/billing/internal/rpc/get_user_balance"
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
	billingCycleRepository := repository.NewBillingCycleRepository(db)
	transactionRepository := repository.NewTransactionRepository(db)

	// Commands and queries
	createTaskCommand := create_task_command.New(taskRepository)
	updateTaskCommand := update_task_command.New(db, taskRepository)
	createUserCommand := create_user_command.New(db, userRepository)
	updateUserCommand := update_user_command.New(db, userRepository)

	createTransactionCommand := create_transaction_command.New(userRepository, billingCycleRepository, transactionRepository)
	createDebitTransactionCommand := create_debit_transaction_command.New(db, taskRepository, createTransactionCommand)
	createCreditTransactionCommand := create_credit_transaction_command.New(db, taskRepository, createTransactionCommand)

	createPayoutTransactionCommand := create_payout_transaction_command.New(db, billingCycleRepository, transactionRepository, createTransactionCommand)

	getUserBalanceQuery := get_user_balance_query.New()

	// Register handlers
	handler := rpc.New(
		get_user_balance.New(getUserBalanceQuery).Handle,
	)

	// Broker consumers
	usersStreamConsumer := users_stream.New(createUserCommand, updateUserCommand)
	go usersStreamConsumer.Consume()

	tasksStreamConsumer := tasks_stream.New(createTaskCommand, updateTaskCommand)
	go tasksStreamConsumer.Consume()

	tasksConsumer := tasks_lifecycle.New(createDebitTransactionCommand, createCreditTransactionCommand)
	go tasksConsumer.Consume()

	// Scheduler
	go runSheduler(createPayoutTransactionCommand)

	// gRPC server
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(middleware.Log(), middleware.JWT(jwtManager)),
	)
	billing.RegisterBillingServer(srv, handler)
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
	err = billing.RegisterBillingHandler(context.Background(), mux, conn)
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

func runSheduler(createPayoutTransactionCommand *create_payout_transaction_command.Command) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return
	}

	defer func() {
		_ = s.Shutdown()
	}()

	// Add a job to the scheduler
	_, err = s.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(23, 59, 59))),
		gocron.NewTask(
			func() {
				//createPayoutTransactionCommand.Handle(context.Background(), nil)
			},
		),
	)
	if err != nil {
		return
	}

	log.Println("scheduler is started")

	s.Start()

	// Wait for a signal to exit
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}
