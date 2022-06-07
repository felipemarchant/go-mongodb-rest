package database

import (
	"context"
	"fmt"
	"go-mongo-rest/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

type MongoClient struct {
	DB *mongo.Client
}

// newInstance from database
func newInstance() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://dev:strongpassword@localhost:27017"))

	if err != nil {
		utils.LogFatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		utils.LogFatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		utils.LogPrintln("MongoDB: Problema ao se conectar")
		return nil
	}

	fmt.Println("MongoDB: Conexão realizado com sucesso")
	return client
}

// Client instance of MongoClient
var Client = &MongoClient{
	DB: newInstance(),
}

// NewInstance from database if you aren't use singleton,
// see Disconnect too.
func (c *MongoClient) NewInstance() *mongo.Client {
	c.DB = newInstance()
	return c.DB
}

// Disconnect from database if you aren't use singleton.
func (c *MongoClient) Disconnect(ctx context.Context) {
	if err := c.DB.Disconnect(ctx); err != nil {
		panic(err)
	}
}

// Ping If you wish to know if a MongoDB server has been found and connected.
func (c *MongoClient) Ping(ctx context.Context) error {
	return c.DB.Ping(ctx, readpref.Primary())
}
