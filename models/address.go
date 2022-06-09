package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Address struct {
	Id      primitive.ObjectID `json:"id" bson:"_id"`
	House   *string            `json:"house_name" bson:"house_name"`
	Street  *string            `json:"street_name" bson:"street_name"`
	City    *string            `json:"city_name" bson:"city_name"`
	PinCode *string            `json:"pin_code" bson:"pin_code"`
}
