package database

import "go.mongodb.org/mongo-driver/mongo"

func (c *MongoClient) UserCollection() *mongo.Collection {
	return c.DB.Database("ecommerce").Collection("users")
}
