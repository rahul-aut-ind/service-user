# service-user

### project setup

#### Prerequisites

- Go 1.20+
- Docker
- MySQL

download dependencies `go mod download`

create a `.env` file taking `.env.example` as reference.

run the service `wire ./... && go run cmd/service-user/main.go`

Quick test from terminal to check service response `curl "localhost:8080/api/v1/users/4"` should yield result after local mysql DB setup and seeding initial data steps listed below is completed.

### local mysql DB setup

pull latest image `docker pull mysql/mysql-server`
run container `docker run --name=mysql-server -d -p 3306:3306 mysql/mysql-server:latest`

check the password from docker logs `docker logs mysql-server`
access mysql cli `docker exec -it mysql-server mysql -u root -p`
create database `CREATE database userdb CHARACTER SET latin1 COLLATE latin1_general_ci;`

create a user and grant permission
`CREATE USER 'root'@'%' IDENTIFIED BY 'some_pass';`
`GRANT ALL PRIVILEGES ON *.* TO 'root'@'%';`

### seeding initial data to mysql DB

copy seeding data to docker container
`docker cp infrastructure/migrations/migration.sql <containerID>:/mysql.sql`
seed data in mysql `source mysql.sql`

### test the service

CREATE USER `curl -X POST "localhost:8080/api/v1/users" -d '{"name":"User14","email":"user14@example.com"}'`
UPDATE USER `curl -X PUT "localhost:8080/api/v1/users/{ID}" -d '{"name":"User14_updated name"'`
GET ALL USERS `curl "localhost:8080/api/v1/users"`
GET SINGLE USER `curl "localhost:8080/api/v1/users/{ID}"`
DELETE USER `curl -X DELETE "localhost:8080/api/v1/users/{ID}"`
