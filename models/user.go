package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	Id           primitive.ObjectID `json:"id" bson:"_id"`
	FirstName    *string            `json:"first_name" validate:"required,min=2,max=30"`
	LastName     *string            `json:"last_name"  validate:"required,min=2,max=30"`
	Password     *string            `json:"password"   validate:"required,min=6"`
	Email        *string            `json:"email"      validate:"email,required"`
	Phone        *string            `json:"phone"      validate:"required"`
	Token        *string            `json:"token"`
	RefreshToken *string            `json:"refresh_token"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	UserCart     []ProductUser      `json:"user_cart" bson:"user_cart"`
	Addresses    []Address          `json:"addresses" bson:"addresses"`
	Orders       []Order            `json:"orders" bson:"orders"`
}
