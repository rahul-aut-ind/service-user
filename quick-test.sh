#!/bin/bash

read -p "Please enter a random number between 500 & 9999: " num

curl "localhost:8080/api/v1/users" -H "x-id-token:something"
response=$(curl -X POST "localhost:8080/api/v1/users" -d '{"firstName":"TestFirstName","lastName":"TestLastName","email":"TestUser'$num'@example.com", "age":19,"address":"Somewhere 10001"}' -H "x-id-token:something")
echo "$response"
id=$(echo "$response" | jq -r '.data.ID')
echo "New user ID is: $id"
curl "localhost:8080/api/v1/users" -H "x-id-token:something"
curl -X PUT "localhost:8080/api/v1/users/$id" -d '{"firstName":"UpadtedFirstName","lastName":"UpdatedLastName", "email":"random@test.com", "age":19,"address":"Str 2, building 5, Floor 9, Flat 10, Somewhere 10001"}' -H "x-id-token:something"
curl "localhost:8080/api/v1/users/$id" -H "x-id-token:something"
curl "localhost:8080/api/v1/users/$id" -H "x-id-token:something"
curl "localhost:8080/api/v1/users/$id" -H "x-id-token:something"
curl -X DELETE "localhost:8080/api/v1/users/$id" -H "x-id-token:something"
curl "localhost:8080/api/v1/users" -H "x-id-token:something"
curl "localhost:8080/api/v1/users/$id" -H "x-id-token:something"