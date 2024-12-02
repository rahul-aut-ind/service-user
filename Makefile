build: sanitize
	docker build -t service-user .

run-service: build
	docker run -d --name=user-service --env-file .env -p 8080:8080 service-user

lint:
	golangci-lint run

test: deps
	go test ./...

test-coverage:
	go test ./... -coverprofile=coverage.out
	sleep 1
	go tool cover -html=coverage.out -o coverage.html

deps:
	wire ./...

sanitize: deps lint test

local-setup:
	docker run --name=mysql-server -d -p 3306:3306 mysql/mysql-server:latest
	echo "mysql server initialized"
	docker run --name=redis-stack -d -p 6379:6379 -p 8001:8001 redis/redis-stack:latest
	echo "redis initialized"
	docker run --name=aws-localstack -d -p 4566:4566 -p 4571:4571 localstack/localstack
	echo "localstack initialized"


local-docker-up:
	docker start mysql-server redis-stack aws-localstack
	echo "waiting for mysql, redis, S3 & dynamoDB to be ready..."

local-docker-down:
	docker stop mysql-server redis-stack aws-localstack

local-docker-delete-service:
	docker stop user-service
	docker rm user-service

local-redeploy: local-docker-delete-service run-service

local-deploy: local-docker-up run-service

service-logs:
	docker logs user-service -f --tail 100

localstack-up:
	docker run --name=aws-localstack -d -p 4566:4566 -p 4571:4571 localstack/localstack
	echo "localstack initialized"

local-aws-setup: localstack-up local-aws-configure local-dynamo-setup local-s3-setup

local-dynamo-setup:
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

local-s3-setup:
	aws --endpoint-url=http://localhost:4566 s3 mb s3://user-images

local-aws-configure:
	aws configure set aws_access_key_id admin
	sleep 1s
	aws configure set aws_secret_access_key password
	sleep 1s
	aws configure set region eu-central-1
	sleep 1s

debug-build-local:
	GOOS=linux GOARCH=arm64/v8 cd ./cmd/service-user && go build -a -o ../../dist/service
	cd ../..

debug-run: lint debug-build-local
	./dist/service

debug-service: deps lint
	go run cmd/service-user/main.go
