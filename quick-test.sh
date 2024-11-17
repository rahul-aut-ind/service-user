#!/bin/bash

read -p "Please enter a random number between 500 & 9999 for userID: " num

echo "-------------"
echo "\nUser Tests"
echo "\nGETALL response below"
curl "localhost:8080/api/v1/users" -H "x-id-token:something"
echo "\n adding an user\n"
response=$(curl -X POST "localhost:8080/api/v1/users" -d '{"firstName":"TestFirstName","lastName":"TestLastName","email":"TestUser'$num'@example.com", "age":19,"address":"Somewhere 10001"}' -H "x-id-token:something")
echo "\nPOST response below"
echo "$response"
id=$(echo "$response" | jq -r '.data.id')
echo "\nNew user ID is: $id"
echo "\nGETALL response below"
curl "localhost:8080/api/v1/users" -H "x-id-token:something"
echo "\nPUT response below"
curl -X PUT "localhost:8080/api/v1/users/$id" -d '{"firstName":"UpadtedFirstName","lastName":"UpdatedLastName", "email":"random@test.com", "age":19,"address":"Str 2, building 5, Floor 9, Flat 10, Somewhere 10001"}' -H "x-id-token:something"
echo "\nGET response below"
curl "localhost:8080/api/v1/users/$id" -H "x-id-token:something"
echo "\nGET response below"
curl "localhost:8080/api/v1/users/$id" -H "x-id-token:something"
echo "\nGET response below"
curl "localhost:8080/api/v1/users/$id" -H "x-id-token:something"
echo "\nDELETE response below"
curl -X DELETE "localhost:8080/api/v1/users/$id" -H "x-id-token:something"
echo "\nGETALL response below"
curl "localhost:8080/api/v1/users" -H "x-id-token:something"
echo "\nGET response below"
curl "localhost:8080/api/v1/users/$id" -H "x-id-token:something"
echo "\n-------------\n"
echo "\nUser Image Tests"
echo "\nGETALL response below"
curl "localhost:8080/api/v1/user-image" -H "x-id-token:something" -H "x-user-id: 11"
echo "\n adding an user image\n"
response=$(curl -X POST 'localhost:8080/api/v1/user-image' \
           --header 'x-user-id: 11' \
           --header 'x-id-token: something' \
           --form 'metadata="{\"takenAt\": \"2024-11-12T00:00:00Z\"}"' \
           --form 'image=@"/Users/rahulupadhyay/Downloads/coins.jpg"')
echo "\nPOST response below"
echo "$response"
id=$(echo "$response" | jq -r '.id')
echo "\nNew image ID is: $id"
echo "\nGET response below"
curl "localhost:8080/api/v1/user-image/$id" -H "x-id-token:something" -H "x-user-id: 11"
echo "\nS3 bucket response"
aws --endpoint-url=http://localhost:4566 s3 ls s3://user-images/story-images/11/
echo "\nGETALL response below"
curl "localhost:8080/api/v1/user-image" -H "x-id-token:something" -H "x-user-id: 11"
echo "\nDELETE response below"
curl -X DELETE "localhost:8080/api/v1/user-image/$id" -H "x-id-token:something" -H "x-user-id: 11"
echo "\nGET response below"
curl "localhost:8080/api/v1/user-image/$id" -H "x-id-token:something" -H "x-user-id: 11"
echo "\nDELETEALL response below"
curl -X DELETE "localhost:8080/api/v1/user-image" -H "x-id-token:something" -H "x-user-id: 11"
echo "\nGET response below"
curl "localhost:8080/api/v1/user-image/$id" -H "x-id-token:something" -H "x-user-id: 11"
echo "\nS3 bucket response"
aws --endpoint-url=http://localhost:4566 s3 ls s3://user-images/story-images/11/
echo "\nDynamo response"
aws dynamodb query --table-name user-images --endpoint-url=http://localhost:4566 \
    --key-condition-expression "UserID = :v1" \
    --expression-attribute-values "{ \":v1\" : { \"S\" : \"11\" } }"