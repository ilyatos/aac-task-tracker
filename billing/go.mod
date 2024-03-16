module github.com/ilyatos/aac-task-tracker/billing

go 1.21

require (
	github.com/go-co-op/gocron/v2 v2.2.5
	github.com/golang-jwt/jwt/v5 v5.2.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1
	github.com/ilyatos/aac-task-tracker/schema_registry v0.0.0-00010101000000-000000000000
	github.com/jackc/pgx/v4 v4.18.1
	github.com/jmoiron/sqlx v1.3.5
	github.com/joho/godotenv v1.5.1
	github.com/satori/go.uuid v1.2.0
	github.com/segmentio/kafka-go v0.4.47
	google.golang.org/genproto/googleapis/api v0.0.0-20240228224816-df926f6c8641
	google.golang.org/grpc v1.62.0
	google.golang.org/protobuf v1.33.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.2 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/jonboulle/clockwork v0.4.0 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/pierrec/lz4/v4 v4.1.15 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/exp v0.0.0-20231219180239-dc181d75b848 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240228201840-1f18d85a4ec2 // indirect
)

replace github.com/ilyatos/aac-task-tracker/schema_registry => /Users/IlyaGoryachev/go/src/aac-task-tracker/schema_registry
