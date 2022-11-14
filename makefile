gen-proto:
	@echo "Generate grpc service from proto files"
	@protoc --go_out=gen/service --go_opt=paths=source_relative --go-grpc_out=gen/service --go-grpc_opt=paths=source_relative --proto_path=proto proto/*.proto

cloud-start:
	docker-compose up -d
cloud-stop:
	docker-compose down -v
cloud-mb: cloud-start
	@sleep 5
	@aws --endpoint-url=http://localhost:4566 s3 mb s3://keeper
cloud-ls:
	@aws --endpoint-url=http://localhost:4566 s3 ls s3://keeper