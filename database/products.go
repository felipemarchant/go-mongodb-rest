package database

import "go.mongodb.org/mongo-driver/mongo"

func (c *MongoClient) ProductCollection() *mongo.Collection {
	return c.DB.Database("ecommerce").Collection("products")
}
