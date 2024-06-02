package repo

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"followers.xws.com/model"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type PersonRepo struct {
	driver neo4j.DriverWithContext
	logger *log.Logger
}

func New(logger *log.Logger) (*PersonRepo, error) {
	uri := os.Getenv("NEO4J_DB")
	user := os.Getenv("NEO4J_USERNAME")
	pass := os.Getenv("NEO4J_PASS")
	auth := neo4j.BasicAuth(user, pass, "")
	// uri := "bolt://localhost:7687"
	// user := "neo4j"
	// pass := "followers"
	// auth := neo4j.BasicAuth(user, pass, "")

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
				userId, _ := record.Get("userId")
				name, _ := record.Get("name")
				surname, _ := record.Get("surname")
				picture, _ := record.Get("picture")
				bio, _ := record.Get("bio")
				quote, _ := record.Get("quote")
				email, _ := record.Get("email")

				person := &model.Person{
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

func (repo *PersonRepo) Follow(userIdToFollow, userIdFollower string) error {
	ctx := context.Background()
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	toFollow, err := strconv.ParseInt(userIdToFollow, 10, 64)
	if err != nil {
		return err
	}
	follower, err := strconv.ParseInt(userIdFollower, 10, 64)
	if err != nil {
		return err
	}

	query := `
        MATCH (p:Person {userId: $followerId}), (p1:Person {userId: $toFollowId})
        MERGE (p)-[r:IS_FOLLOWING]->(p1)
        ON CREATE SET r.since = $date
        RETURN r
    `

	_, err = session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				query,
				map[string]interface{}{"followerId": follower, "toFollowId": toFollow, "date": time.Now().Format("2006-01-02")})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})
	if err != nil {
		repo.logger.Println("Error inserting relationship:", err)
		return err
	}
	repo.logger.Println("Relationship created successfully")
	return nil
}

func (repo *PersonRepo) UnFollow(userIdToUnFollow, userIdFollower string) error {
	ctx := context.Background()
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	toUnFollow, err := strconv.ParseInt(userIdToUnFollow, 10, 64)
	if err != nil {
		return err
	}
	follower, err := strconv.ParseInt(userIdFollower, 10, 64)
	if err != nil {
		return err
	}

	query := `
		MATCH (p:Person {userId: $followerId})-[r:IS_FOLLOWING]->(p1:Person {userId: $toUnFollowId})
		DELETE r
	`

	_, err = session.ExecuteWrite(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				query,
				map[string]interface{}{"followerId": follower, "toUnFollowId": toUnFollow})
			if err != nil {
				return nil, err
			}

			if result.Next(ctx) {
				return result.Record().Values[0], nil
			}

			return nil, result.Err()
		})
	if err != nil {
		repo.logger.Println("Error unfollowing:", err)
		return err
	}
	repo.logger.Println("Relationship deleted successfully")
	return nil
}

func (repo *PersonRepo) GetFollowers(userId string) (model.Followers, error) {
	ctx := context.Background()
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	userIdInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return nil, err
	}

	query := `
		MATCH (f:Person)-[:IS_FOLLOWING]->(p:Person {userId: $userId})
		RETURN f.userId as userId, f.name as name, f.surname as surname, f.quote as quote, f.email as email
	`

	followersResults, err := session.ExecuteRead(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				query,
				map[string]interface{}{"userId": userIdInt})
			if err != nil {
				return nil, err
			}
			var followers model.Followers
			for result.Next(ctx) {
				record := result.Record()
				userId, _ := record.Get("userId")
				name, _ := record.Get("name")
				surname, _ := record.Get("surname")
				quote, _ := record.Get("quote")
				email, _ := record.Get("email")

				person := &model.Follower{
					UserId:  userId.(int64),
					Name:    name.(string),
					Surname: surname.(string),
					Quote:   quote.(string),
					Email:   email.(string),
				}
				followers = append(followers, person)
			}
			return followers, nil
		})
	if err != nil {
		repo.logger.Println("Error querying search:", err)
		return nil, err
	}
	return followersResults.(model.Followers), nil
}

func (repo *PersonRepo) GetRecommended(userId string) (model.Followers, error) {
	ctx := context.Background()
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	userIdInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return nil, err
	}

	query := `
		MATCH (user:Person {userId: $userId})-[:IS_FOLLOWING]->(following:Person)-[:IS_FOLLOWING]->(f:Person)
		WHERE NOT (user)-[:IS_FOLLOWING]->(f)
	  	AND user <> following
		RETURN f.userId as userId, f.name as name, f.surname as surname, f.quote as quote, f.email as email
	
	`

	followersResults, err := session.ExecuteRead(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				query,
				map[string]interface{}{"userId": userIdInt})
			if err != nil {
				return nil, err
			}
			var followers model.Followers
			for result.Next(ctx) {
				record := result.Record()
				userId, _ := record.Get("userId")
				name, _ := record.Get("name")
				surname, _ := record.Get("surname")
				quote, _ := record.Get("quote")
				email, _ := record.Get("email")

				person := &model.Follower{
					UserId:  userId.(int64),
					Name:    name.(string),
					Surname: surname.(string),
					Quote:   quote.(string),
					Email:   email.(string),
				}
				followers = append(followers, person)
			}
			return followers, nil
		})
	if err != nil {
		repo.logger.Println("Error querying search:", err)
		return nil, err
	}
	return followersResults.(model.Followers), nil
}

func (repo *PersonRepo) GetFollowing(userId string) (model.Followers, error) {
	ctx := context.Background()
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	userIdInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return nil, err
	}

	query := `
		MATCH (f:Person {userId: $userId})-[:IS_FOLLOWING]->(p:Person)
		RETURN p.userId as userId, p.name as name, p.surname as surname, p.quote as quote, p.email as email
	`

	followersResults, err := session.ExecuteRead(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				query,
				map[string]interface{}{"userId": userIdInt})
			if err != nil {
				return nil, err
			}
			var followers model.Followers
			for result.Next(ctx) {
				record := result.Record()
				userId, _ := record.Get("userId")
				name, _ := record.Get("name")
				surname, _ := record.Get("surname")
				quote, _ := record.Get("quote")
				email, _ := record.Get("email")

				person := &model.Follower{
					UserId:  userId.(int64),
					Name:    name.(string),
					Surname: surname.(string),
					Quote:   quote.(string),
					Email:   email.(string),
				}
				followers = append(followers, person)
			}
			return followers, nil
		})
	if err != nil {
		repo.logger.Println("Error querying search:", err)
		return nil, err
	}
	return followersResults.(model.Followers), nil
}

func (repo *PersonRepo) IsFollowing(userId, followingUserId string) (bool, error) {
	ctx := context.Background()
	session := repo.driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: "neo4j"})
	defer session.Close(ctx)

	userIdInt, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		return false, err
	}

	followingUserIdInt, err := strconv.ParseInt(followingUserId, 10, 64)
	if err != nil {
		return false, err
	}

	query := `
		MATCH (f:Person {userId: $userId})-[:IS_FOLLOWING]->(p:Person {userId: $followingUserId})
		RETURN p.name as name
	`

	followersResults, err := session.ExecuteRead(ctx,
		func(transaction neo4j.ManagedTransaction) (any, error) {
			result, err := transaction.Run(ctx,
				query,
				map[string]interface{}{"userId": userIdInt, "followingUserId": followingUserIdInt})
			if err != nil {
				return false, err
			}
			if result.Next(ctx) {
				return true, nil
			}
			return false, nil
		})
	if err != nil {
		repo.logger.Println("Error querying search:", err)
		return false, err
	}
	return followersResults.(bool), nil
}
