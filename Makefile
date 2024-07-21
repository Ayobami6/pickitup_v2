.PHONY: gen-users start-users gen-riders start-riders start-gateway gen-orders 

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

start-orders:
	@cd orders && air

gen-riders:
	@protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    common/proto/riders/riders.proto

gen-orders:
	@protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    common/proto/orders/orders.proto