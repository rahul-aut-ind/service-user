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