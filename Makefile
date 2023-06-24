build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./.bin/app ./cmd/main.go

run:
	go run ./cmd/media-service/main.go

test:
	go test -v -count=1 ./...
start: build
	docker-compose up --build media-service

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

lint:
	echo "Starting linters"
	golangci-lint run ./...

migrate-up:
	migrate -database postgres://postgres:postgres@localhost:5432/media?sslmode=disable -path migrations up 1

PHONY: generate
generate:
	mkdir -p pkg/payment-service
	protoc --go_out=./pkg/payment-service --go_opt=paths=source_relative \
	       --go-grpc_out=./pkg/payment-service --go-grpc_opt=paths=source_relative \
		   api/payment-service/product.proto  
	mv pkg/payment-service/api/payment-service/* pkg/payment-service
	rm -rf pkg/payment-service/api
