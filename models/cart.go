package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cart struct {
	Id     primitive.ObjectID `json:"id" bson:"_id"`
	Name   *string            `json:"name" bson:"name"`
	Price  int                `json:"price"  bson:"price"`
	Rating *uint              `json:"rating" bson:"rating"`
	Image  *string            `json:"image"  bson:"image"`
}
