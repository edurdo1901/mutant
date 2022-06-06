package utilconection

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// OpenDB create database connection.
func OpenDB(dbconection, dataBase string) (*mongo.Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(dbconection))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)

	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		return nil, err
	}

	return client.Database(dataBase), nil
}
