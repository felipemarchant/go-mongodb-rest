package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Order struct {
	OrderId       primitive.ObjectID `bson:"_id"`
	OrderCart     []ProductUser      `json:"order_list"  bson:"order_list"`
	OrdereredAt   time.Time          `json:"ordered_on"  bson:"ordered_on"`
	Price         int                `json:"total_price" bson:"total_price"`
	Discount      *int               `json:"discount"    bson:"discount"`
	PaymentMethod Payment            `json:"payment_method" bson:"payment_method"`
}
