package repo

import (
	"context"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type PersonRepo struct {
	driver neo4j.DriverWithContext
	logger *log.Logger
}

func New(logger *log.Logger) (*PersonRepo, error) {
	//TODO: Initialize connection to db
	uri := "bolt://localhost:7687"
	user := "neo4j"
	pass := "followers"
	auth := neo4j.BasicAuth(user, pass, "")

	driver, err := neo4j.NewDriverWithContext(uri, auth)
	if err != nil {
		logger.Panic(err)
		return nil, err
	}

	return &PersonRepo{
		driver: driver,
		logger: logger,
	}, nil
}

func (mr *PersonRepo) CheckConnection() {
	ctx := context.Background()
	err := mr.driver.VerifyConnectivity(ctx)
	if err != nil {
		mr.logger.Panic(err)
		return
	}
	mr.logger.Printf(`Neo4J server address: %s`, mr.driver.Target().Host)
}

func (mr *PersonRepo) CloseDriverConnection(ctx context.Context) {
	mr.driver.Close(ctx)
}
