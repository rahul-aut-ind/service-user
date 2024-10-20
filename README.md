# service-user

### local mysql DB setup

pull latest image `docker pull mysql/mysql-server`
run container `docker run --name=mysql-server -d -p 3306:3306 mysql/mysql-server:latest`

check the password from docker logs `docker logs mysql-server`
access mysql cli `docker exec -it mysql-server mysql -u root -p`
create database `CREATE database userdb CHARACTER SET latin1 COLLATE latin1_general_ci;`

create a user and grant permission
`CREATE USER 'root'@'%' IDENTIFIED BY 'some_pass';`
`GRANT ALL PRIVILEGES ON *.* TO 'root'@'%';`

### seeding data to mysql DB

copy seeding data to docker container
`docker cp infrastructure/migrations/migration.sql <containerID>:/mysql.sql`
seed data in mysql `source mysql.sql`
