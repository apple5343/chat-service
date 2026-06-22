generate-auth-api:
	protoc --proto_path api/auth_v1 --go_out=pkg/api/auth_v1 --go_opt=paths=source_relative --go-grpc_out=pkg/api/auth_v1 --go-grpc_opt=paths=source_relative api/auth_v1/auth.proto

generate-user-api:
	protoc --proto_path api/user_v1 --go_out=pkg/api/user_v1 --go_opt=paths=source_relative --go-grpc_out=pkg/api/user_v1 --go-grpc_opt=paths=source_relative api/user_v1/user.proto
