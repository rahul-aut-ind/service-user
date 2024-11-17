build: sanitize
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

sanitize: deps lint test

local-docker-up:
	docker start mysql-server redis-stack
	echo "waiting for mysql and redis to be ready..."
	sleep 1s
	docker start user-service

local-docker-down:
	docker stop mysql-server redis-stack user-service

local-docker-delete-service:
	docker stop user-service
	docker rm user-service

local-redeploy: local-docker-delete-service run-service local-docker-up

local-clean-run: run-service local-docker-up

service-logs:
	docker logs user-service -f --tail 100

localstack-up:
	docker run --rm --name=localstack -p 4566:4566 -p 4571:4571 localstack/localstack

local-dynamo-setup: localstack-up local-aws-setup
	aws --endpoint-url=http://localhost:4566 dynamodb create-table \
        --table-name user-images \
        --attribute-definitions \
            AttributeName=UserID,AttributeType=S \
            AttributeName=ImageID,AttributeType=S \
            AttributeName=TakenAt,AttributeType=S \
        --key-schema \
            AttributeName=UserID,KeyType=HASH \
            AttributeName=ImageID,KeyType=RANGE \
        --provisioned-throughput \
            ReadCapacityUnits=5,WriteCapacityUnits=5 \
        --local-secondary-indexes \
            "[{\"IndexName\": \"UserIDTakenAtIndex\", \"KeySchema\":[{\"AttributeName\":\"UserID\",\"KeyType\":\"HASH\"}, {\"AttributeName\":\"TakenAt\",\"KeyType\":\"RANGE\"}],\"Projection\":{\"ProjectionType\":\"ALL\"}}]" \
        --table-class STANDARD

local-s3-setup: localstack-up local-aws-setup
	aws --endpoint-url=http://localhost:4566 s3 mb s3://user-images

local-aws-setup:
	aws configure set aws_access_key_id admin
	sleep 1s
	aws configure set aws_secret_access_key password
	sleep 1s
	aws configure set region eu-central-1
	sleep 1s

