package integrationtest

import (
	"context"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

type MySQLSetup struct {
	container  *mysql.MySQLContainer
	ConnString string
}

func NewMySQLSetup() *MySQLSetup {
	setup := &MySQLSetup{}
	setup.initialize()
	return setup
}

func (tc *MySQLSetup) initialize() {
	ctx := context.Background()
	tc.createMySQLContainer(ctx)
}

func (tc *MySQLSetup) createMySQLContainer(ctx context.Context) {
	log.Printf("starting test container")
	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0.36",
		mysql.WithDatabase("usersdb"),
		mysql.WithUsername("root"),
		mysql.WithPassword("password"),
	)
	if err != nil {
		log.Printf("failed to start container: %s", err)
		return
	}

	connString, err := mysqlContainer.ConnectionString(ctx, "charset=utf8mb4", "parseTime=True", "loc=Local")
	if err != nil {
		log.Printf("failed to get connection string : %s", err)
		return
	}

	tc.container = mysqlContainer
	tc.ConnString = connString
	log.Printf("test container started")
}

func (tc *MySQLSetup) Stop() {
	if err := testcontainers.TerminateContainer(tc.container); err != nil {
		log.Printf("failed to terminate container: %s", err)
	}
}
