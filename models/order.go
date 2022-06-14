package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Order struct {
	Id            primitive.ObjectID `json:"id" bson:"_id"`
	Cart          []Cart             `json:"cart"  bson:"cart"`
	CreatedAt     time.Time          `json:"created_at"  bson:"created_at"`
	Price         int                `json:"total_price" bson:"total_price"`
	Discount      *int               `json:"discount"    bson:"discount"`
	PaymentMethod Payment            `json:"payment_method" bson:"payment_method"`
}
