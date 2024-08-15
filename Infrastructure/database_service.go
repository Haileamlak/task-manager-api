package infrastructure

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// DatabaseService interface
type DatabaseService interface {
	Connect(dbUri string) *mongo.Database
}

// databaseService struct
type databaseService struct {
}

// NewDatabase creates a new database service
func NewDatabase() DatabaseService {
	return &databaseService{}
}

// Connect connects to the database
func (d *databaseService) Connect(dbUri string) *mongo.Database {
	clientOptions := options.Client().ApplyURI(dbUri)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database("task_manager")
}