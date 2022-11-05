gen-proto:
	@echo "Generate grpc service from proto files"
	@protoc --go_out=gen/service --go_opt=paths=source_relative --go-grpc_out=gen/service --go-grpc_opt=paths=source_relative --proto_path=proto proto/*.proto