module github.com/kozlovsv/evaluator/server

go 1.21.5

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/joho/godotenv v1.5.1
	google.golang.org/grpc v1.63.0
)

require (
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240227224415-6ceb2ff114de // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)

require (
	evaluator/protos v0.0.0-00010101000000-000000000000
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.1.0
)

replace evaluator/protos => ../protos
