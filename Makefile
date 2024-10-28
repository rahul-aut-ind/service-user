build:
	docker build -t service-user .

run-service: build
	docker run -d --name=user-service --env-file .env -p 8080:8080 service-user

build-local:
	GOOS=linux GOARCH=arm64/v8 cd ./cmd/service-user && go build -a -o ../../dist/service
	cd ../..

local-run: lint build-local
	./dist/service

debug-service: deps lint
	go run cmd/service-user/main.go

lint:
	golangci-lint run

test: deps
	go test ./...

deps:
	wire ./...

local-setup:
	docker run --name=mysql-server -d -p 3306:3306 mysql/mysql-server:latest
	echo "mysql server initialized"
	docker run -d --name redis-stack -p 6379:6379 -p 8001:8001 redis/redis-stack:latest
	echo "redis initialized"