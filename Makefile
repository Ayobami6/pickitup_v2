.PHONY: gen-users

gen-users:
	@protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    common/proto/users/users.proto

