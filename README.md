# service-user

### project setup

download dependencies `go mod download`

create a `.env` file taking `.env.example` as reference.

### run the service on docker (locally)

run the service from terminal `make local-deploy`

### unit & integration tests

run in terminal `make test`

### test coverage

run in terminal `make test-coverage` and open the coverage.html file

### lint check

run in terminal `make lint`
make sure to `brew install golangci-lint` before

##### debug the service locally

run the service from terminal `make debug-run`

### prerequisites

- Go 1.23+
- Docker
- MySQL 8.0
- Redis
- AWS S3
- AWS DynamoDB

### local AWS services setup

run in terminal `make local-aws-setup`

### local redis & mysql setup

run in terminal `make local-setup`

##### check the password from docker logs

`docker logs mysql-server`

##### access mysql cli

`docker exec -it mysql-server mysql -u root -p`  || (reset pswd for first time use)

##### create database

`CREATE database userdb CHARACTER SET latin1 COLLATE latin1_general_ci;`

##### create a user and grant permission

`CREATE USER 'root'@'%' IDENTIFIED BY 'some_pass';`

`GRANT ALL PRIVILEGES ON *.* TO 'root'@'%';`

### seeding initial data to mysql DB

##### copy seeding data to docker container

`docker cp infrastructure/migrations/migration.sql <containerID>:/mysql.sql`

##### seed data in mysql

`source mysql.sql` from mysql prompt

---

### test the service

To test indivizual endpoints refer below:
```sh
##### CREATE USER

`curl -X POST "localhost:8080/api/v1/users" -d '{"firstName":"TestFirstName","lastName":"TestLastName","email":"TestUser'$num'@example.com", "age":19,"address":"Somewhere 10001"}' -H "x-id-token:something"`
```
```sh
##### UPDATE USER

`curl -X PUT "localhost:8080/api/v1/users/$id" -d '{"firstName":"UpadtedFirstName","lastName":"UpdatedLastName", "email":"random@test.com", "age":19,"address":"Str 2, building 5, Floor 9, Flat 10, Somewhere 10001"}' -H "x-id-token:something"`
```
```sh
##### GET ALL USERS

`curl "localhost:8080/api/v1/users" -H "x-id-token:something"`
```
```sh
##### GET SINGLE USER

`curl "localhost:8080/api/v1/users/$id" -H "x-id-token:something"`
```
```sh
##### DELETE USER

`curl -X DELETE "localhost:8080/api/v1/users/$id" -H "x-id-token:something"`
```
```sh
##### CREATE USER IMAGE

`curl -X POST 'localhost:8080/api/v1/user-image' \
--header 'x-user-id: 11' \
--header 'x-id-token: something' \
--form 'metadata="{\"takenAt\": \"2024-11-12T00:00:00Z\"}"' \
--form 'image=@"/Users/rahulupadhyay/Downloads/coins.jpg"'`
```
```sh
##### GET ALL USER IMAGES

`curl "localhost:8080/api/v1/user-image" -H "x-id-token:something" -H "x-user-id: 11"`
```
```sh
##### GET SINGLE USER IMAGE

`curl "localhost:8080/api/v1/user-image/$id" -H "x-id-token:something" -H "x-user-id: 11"`
```
```sh
##### DELETE USER IMAGE

`curl -X DELETE "localhost:8080/api/v1/user-image/$id" -H "x-id-token:something" -H "x-user-id: 11"`
```
```sh
##### DELETE ALL USER IMAGES

`curl -X DELETE "localhost:8080/api/v1/user-image" -H "x-id-token:something" -H "x-user-id: 11"`
```


###### Note: 
- For a quick test of all apis, open a terminal window and run `sh quick-test.sh`

----