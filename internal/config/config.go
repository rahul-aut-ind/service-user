package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type (
	Env struct {
		// Environment the development environment
		Environment string
		// DBConnectionString the connections string
		DBConnectionString string
		// ServerHost the host that the server will start on
		ServerHost string
		// ServerPort the port that server will start on
		ServerPort string
		// RedisAddress the server and port where redis will run
		RedisAddress string
		// DefaultAWSRegion default AWS region
		DefaultAWSRegion string
		// DynamoDBConnectionString the connections string
		DynamoDBConnectionString string
		// DynamoDBTable is the table name in dynamoDB
		DynamoDBTable string
		// AwsAccessKey is the aws access key
		AwsAccessKey string
		// AwsSecretAccessKey is the Aws SecretAccess Key
		AwsSecretAccessKey string
		// S3Bucket is the S3 bucket name
		S3Bucket string
		// S3Directory is the S3 directory in the bucket
		S3Directory string
	}
)

const (
	// LocalEnvironment is the local dev environment
	LocalEnvironment = "development"
	// HeaderUserID name of the header that holds the user id
	HeaderUserID = "x-user-id"
	// HeaderIDToken name of the header that holds the id token
	HeaderIDToken = "x-id-token"
	// HeaderContentType name of the header that holds the content type
	HeaderContentType = "content-type"
	// QueryParamLastKey name of query param that holds last evaluated key
	QueryParamLastKey = "lastKey"
	// QueryParamlastKeyDate name of query param that holds last evaluated key date
	QueryParamlastKeyDate = "lastKeyDate"
	// QueryParamLimit name of query param that holds history limit
	QueryParamLimit = "limit"
)

// NewEnv creates a new instance of Env
// tries to load the env variables from .env
func NewEnv() *Env {
	path, err := os.Getwd()
	if err != nil {
		log.Fatalf("error getting path")
	}
	dotenvError := godotenv.Load(fmt.Sprintf("%s/.env", path))
	if dotenvError != nil {
		log.Printf("error loading .env file, ignoring dotenv")
	}

	return &Env{
		DBConnectionString:       os.Getenv("MysqlDB_Connection_String"),
		ServerHost:               os.Getenv("Server_Host"),
		ServerPort:               os.Getenv("Server_Port"),
		RedisAddress:             os.Getenv("Redis_Address"),
		DefaultAWSRegion:         os.Getenv("AWS_REGION"),
		DynamoDBConnectionString: os.Getenv("DynamoDB_Connection_String"),
		Environment:              os.Getenv("Environment"),
		DynamoDBTable:            os.Getenv("DynamoDB_Table"),
		AwsAccessKey:             os.Getenv("AWS_ACCESS_KEY_ID"),
		AwsSecretAccessKey:       os.Getenv("AWS_SECRET_ACCESS_KEY"),
		S3Bucket:                 os.Getenv("S3Bucket"),
		S3Directory:              os.Getenv("S3Directory"),
	}
}
