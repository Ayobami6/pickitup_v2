.PHONY: gen-users start-users

start-users:
	@cd users && air

gen-users:
	@protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    common/proto/users/users.proto

start-gateway:
	@cd gateway && air

start-riders:
	@cd riders && air

gen-riders:
	@protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    common/proto/riders/riders.proto