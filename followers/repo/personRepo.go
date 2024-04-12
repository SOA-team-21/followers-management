package repo

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"followers.xws.com/model"
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

func (pr *PersonRepo) CheckConnection() {
	ctx := context.Background()
	err := pr.driver.VerifyConnectivity(ctx)
	if err != nil {
		pr.logger.Panic(err)
		return
	}
	pr.logger.Printf(`Neo4J server address: %s`, pr.driver.Target().Host)
}

func (pr *PersonRepo) CloseDriverConnection(ctx context.Context) {
	pr.driver.Close(ctx)
}

func (pr *PersonRepo) GetPerson(userId string) (*model.Person, error) {
	// Neo4J Sessions are lightweight so we create one for each transaction
	// Sessions are NOT thread safe
	ctx := context.Background()
	session := pr.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	userIdInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return nil, err
	}

	query := `
		MATCH (p:Person {userId: $userId})
		RETURN p.id as id, p.userId as userId, p.name as name, p.surname as surname,
		p.picture as picture, p.bio as bio, p.quote as quote, p.email as email
		LIMIT 1
	`

	personResult, err := session.ExecuteRead(ctx,
		func(transaction neo4j.ManagedTransaction) (interface{}, error) {
			result, err := transaction.Run(ctx, query, map[string]interface{}{"userId": userIdInt})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				record := result.Record()
				id, _ := record.Get("id")
				userId, _ := record.Get("userId")
				name, _ := record.Get("name")
				surname, _ := record.Get("surname")
				picture, _ := record.Get("picture")
				bio, _ := record.Get("bio")
				quote, _ := record.Get("quote")
				email, _ := record.Get("email")

				person := &model.Person{
					Id:      id.(int64),
					UserId:  userId.(int64),
					Name:    name.(string),
					Surname: surname.(string),
					Picture: picture.(string),
					Bio:     bio.(string),
					Quote:   quote.(string),
					Email:   email.(string),
				}
				return person, nil
			}
			return nil, err
		})

	if err != nil {
		pr.logger.Println("Error querying Person:", err)
		return nil, err
	}
	if person, ok := personResult.(*model.Person); ok {
		return person, nil
	}
	return nil, fmt.Errorf("unexpected result type from Neo4j query")
}