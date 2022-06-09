package controllers

import (
	"context"
	"github.com/felipemarchant/go-mongo-rest/database"
	"github.com/felipemarchant/go-mongo-rest/models"
	r "github.com/felipemarchant/go-mongo-rest/rest"
	"github.com/felipemarchant/go-mongo-rest/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

var msgInvalidCode = "Código inválido"

func AddAddress(c *gin.Context) {
	userId := c.Query("id")

	if userId == "" {
		r.Response(c, msgInvalidCode, http.StatusNotFound)
		return
	}

	address, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		r.Response(c, err.Error(), http.StatusInternalServerError)
		return
	}

	var addresses models.Address
	addresses.Id = primitive.NewObjectID()

	if err = c.BindJSON(&addresses); err != nil {
		r.Response(c, err.Error(), http.StatusNotAcceptable)
		return
	}

	var ctx, cancel = utils.ContextWithTimeout()
	defer cancel()
	defer ctx.Done()

	users := database.Client.UserCollection()

	matchFilter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: address}}}}
	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$addresses"}}}}
	group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "count", Value: bson.D{primitive.E{Key: "$sum", Value: 1}}}}}}

	cursor, err := users.Aggregate(ctx, mongo.Pipeline{matchFilter, unwind, group})
	if err != nil {
		r.Response(c, err.Error(), http.StatusInternalServerError)
		return
	}

	var addressInfo []bson.M
	if err = cursor.All(ctx, &addressInfo); err != nil {
		r.Response(c, err.Error(), http.StatusInternalServerError)
		return
	}

	var size int32
	for _, addressNo := range addressInfo {
		count := addressNo["count"]
		size = count.(int32)
	}

	if size < 2 {
		filter := bson.D{primitive.E{Key: "_id", Value: address}}
		update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "addresses", Value: addresses}}}}
		_, err := users.UpdateOne(ctx, filter, update)

		if err != nil {
			r.Response(c, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		r.Response(c, "Limite atingido", http.StatusBadRequest)
		return
	}

	r.Response(c, addresses, http.StatusCreated)
}

func EditHomeAddress(c *gin.Context) {
	userId := c.Query("id")

	if userId == "" {
		r.Response(c, msgInvalidCode, http.StatusNotFound)
		return
	}

	userIdx, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		r.Response(c, err.Error(), http.StatusInternalServerError)
		return
	}

	var editAddress models.Address
	if err := c.BindJSON(&editAddress); err != nil {
		r.Response(c, err.Error(), http.StatusBadRequest)
		return
	}

	var ctx, cancel = utils.ContextWithTimeout()
	defer cancel()
	defer ctx.Done()

	users := database.Client.UserCollection()

	filter := bson.D{primitive.E{Key: "_id", Value: userIdx}}
	update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "addresses.0.house_name", Value: editAddress.House}, {Key: "addresses.0.street_name", Value: editAddress.Street}, {Key: "addresses.0.city_name", Value: editAddress.City}, {Key: "addresses.0.pin_code", Value: editAddress.PinCode}}}}

	_, err = users.UpdateOne(ctx, filter, update)
	if err != nil {
		r.Response(c, err.Error(), http.StatusInternalServerError)
		return
	}

	r.Response(c, "Endereço de casa atualizado com sucesso", http.StatusOK)
}

func EditWorkAddress(c *gin.Context) {
	userId := c.Query("id")

	if userId == "" {
		r.Response(c, msgInvalidCode, http.StatusNotFound)
		return
	}

	userIdx, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		r.Response(c, err.Error(), http.StatusInternalServerError)
		return
	}

	var editAddress models.Address
	if err := c.BindJSON(&editAddress); err != nil {
		r.Response(c, err.Error(), http.StatusBadRequest)
		return
	}

	var ctx, cancel = utils.ContextWithTimeout()
	defer cancel()
	defer ctx.Done()
	filter := bson.D{primitive.E{Key: "_id", Value: userIdx}}
	update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "addresses.1.house_name", Value: editAddress.House}, {Key: "addresses.1.street_name", Value: editAddress.Street}, {Key: "addresses.1.city_name", Value: editAddress.City}, {Key: "addresses.1.pin_code", Value: editAddress.PinCode}}}}

	users := database.Client.UserCollection()

	_, err = users.UpdateOne(ctx, filter, update)
	if err != nil {
		r.Response(c, err.Error(), http.StatusInternalServerError)
		return
	}

	r.Response(c, "Endereço de trabalho atualizado com sucesso", http.StatusOK)
}

func DeleteAddress(c *gin.Context) {
	userId := c.Query("id")

	if userId == "" {
		r.Response(c, msgInvalidCode, http.StatusNotFound)
		return
	}

	addresses := make([]models.Address, 0)
	userIdx, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		r.Response(c, err.Error(), http.StatusInternalServerError)
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	defer ctx.Done()

	filter := bson.D{primitive.E{Key: "_id", Value: userIdx}}
	update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "addresses", Value: addresses}}}}

	users := database.Client.UserCollection()

	_, err = users.UpdateOne(ctx, filter, update)
	if err != nil {
		r.Response(c, err.Error(), http.StatusNotFound)
		return
	}

	r.Response(c, "Endereço removido com sucesso", http.StatusOK)
}
