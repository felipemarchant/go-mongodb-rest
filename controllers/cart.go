package controllers

import (
	"context"
	"github.com/felipemarchant/go-mongo-rest/database"
	"github.com/felipemarchant/go-mongo-rest/models"
	r "github.com/felipemarchant/go-mongo-rest/rest"
	"github.com/felipemarchant/go-mongo-rest/security"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

func AddToCart(c *gin.Context) {
	body := make(map[string]interface{})
	err := c.BindJSON(&body)
	if err != nil {
		r.Response(c, err.Error(), http.StatusBadRequest)
		return
	}

	productId := body["product_id"].(string)
	if productId == "" {
		r.Response(c, "product id is empty", http.StatusBadRequest)
		return
	}

	principal := security.UserPrincipal(c)

	productID, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		r.Response(c, err.Error(), http.StatusInternalServerError)
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = addToCart(ctx, productID, principal.Id)
	if err != nil {
		r.Response(c, err.Error(), http.StatusInternalServerError)
		return
	}

	r.Response(c, "Produto adicionado ao carrinho com sucesso", http.StatusCreated)
}

func addToCart(ctx context.Context, productId primitive.ObjectID, userId primitive.ObjectID) error {
	products := database.Client.ProductCollection()
	users := database.Client.UserCollection()

	productFound, err := products.Find(ctx, bson.M{"_id": productId})
	if err != nil {
		return err
	}

	var productCart []models.Cart
	err = productFound.All(ctx, &productCart)
	if err != nil {
		return err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: userId}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "user_cart", Value: bson.D{{Key: "$each", Value: productCart}}}}}}
	_, err = users.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFromCart(c *gin.Context) {
	productId := c.Param("product")
	if productId == "" {
		r.Response(c, "product id is empty", http.StatusBadRequest)
		return
	}

	principal := security.UserPrincipal(c)

	productID, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		r.Response(c, err.Error(), http.StatusInternalServerError)
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = deleteFromCart(ctx, productID, principal.Id)
	if err != nil {
		r.Response(c, err.Error(), http.StatusInternalServerError)
		return
	}

	r.Response(c, "Produto removido do carrinho com sucesso", http.StatusOK)
}

func deleteFromCart(ctx context.Context, productId primitive.ObjectID, userId primitive.ObjectID) error {
	users := database.Client.UserCollection()
	filter := bson.D{primitive.E{Key: "_id", Value: userId}}
	update := bson.M{"$pull": bson.M{"user_cart": bson.M{"_id": productId}}}
	_, err := users.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func ListCartItem(c *gin.Context) {
	userId := security.UserPrincipal(c).Id

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	users := database.Client.UserCollection()

	var filledCart *models.User
	err := users.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: userId}}).Decode(&filledCart)
	if err != nil {
		r.Response(c, err.Error(), http.StatusInternalServerError)
		return
	}

	status := http.StatusOK
	if len(filledCart.UserCart) < 1 {
		status = http.StatusNotFound
	}

	r.Response(c, filledCart.UserCart, status)
}

func InstantCheckout(c *gin.Context) {
	principal := security.UserPrincipal(c)

	body := make(map[string]interface{})
	err := c.BindJSON(&body)
	if err != nil {
		r.Response(c, err.Error(), http.StatusBadRequest)
		return
	}

	productId := body["product_id"].(string)
	if productId == "" {
		r.Response(c, "product id is empty", http.StatusBadRequest)
		return
	}

	productID, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		r.Response(c, err.Error(), http.StatusInternalServerError)
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = instantCheckout(ctx, productID, principal.Id)
	if err != nil {
		r.Response(c, err.Error(), http.StatusInternalServerError)
		return
	}

	r.Response(c, "Compra realizada com sucesso", http.StatusCreated)
}

func instantCheckout(ctx context.Context, productId primitive.ObjectID, userId primitive.ObjectID) error {
	var productDetails models.Cart
	var orderDetail models.Order
	orderDetail.Id = primitive.NewObjectID()
	orderDetail.CreatedAt = time.Now()
	orderDetail.Cart = make([]models.Cart, 0)
	orderDetail.PaymentMethod.Cod = true

	products := database.Client.ProductCollection()
	users := database.Client.UserCollection()

	err := products.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: productId}}).Decode(&productDetails)
	if err != nil {
		log.Println(err)
	}

	orderDetail.Price = productDetails.Price
	filter := bson.D{primitive.E{Key: "_id", Value: userId}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderDetail}}}}
	_, err = users.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: userId}}
	update2 := bson.M{"$push": bson.M{"orders.$[].cart": productDetails}}
	_, err = users.UpdateOne(ctx, filter2, update2)
	if err != nil {
		return err
	}

	return nil
}

func CheckoutCart(c *gin.Context) {
	userId := security.UserPrincipal(c).Id

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err := checkoutCart(ctx, userId)
	if err != nil {
		r.Response(c, err.Error(), http.StatusInternalServerError)
		return
	}

	r.Response(c, "Compra realizada com sucesso", http.StatusCreated)
}

func checkoutCart(ctx context.Context, userId primitive.ObjectID) error {
	var cartItems models.User
	var orderCart models.Order
	orderCart.Id = primitive.NewObjectID()
	orderCart.CreatedAt = time.Now()
	orderCart.Cart = make([]models.Cart, 0)
	orderCart.PaymentMethod.Cod = true

	users := database.Client.UserCollection()

	unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$user_cart"}}}}
	grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$user_cart.price"}}}}}}
	currentResults, err := users.Aggregate(ctx, mongo.Pipeline{unwind, grouping})
	ctx.Done()
	if err != nil {
		return err
	}

	var userCart []bson.M
	if err = currentResults.All(ctx, &userCart); err != nil {
		return err
	}

	var totalPrice int32
	for _, userItem := range userCart {
		price := userItem["total"]
		totalPrice = price.(int32)
	}

	orderCart.Price = int(totalPrice)
	filter := bson.D{primitive.E{Key: "_id", Value: userId}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: orderCart}}}}
	_, err = users.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	err = users.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: userId}}).Decode(&cartItems)
	if err != nil {
		return err
	}

	filter2 := bson.D{primitive.E{Key: "_id", Value: userId}}
	update2 := bson.M{"$push": bson.M{"orders.$[].cart": bson.M{"$each": cartItems.UserCart}}}
	_, err = users.UpdateOne(ctx, filter2, update2)
	if err != nil {
		return err
	}

	userCartEmpty := make([]models.Cart, 0)
	filtered := bson.D{primitive.E{Key: "_id", Value: userId}}
	updated := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "user_cart", Value: userCartEmpty}}}}
	_, err = users.UpdateOne(ctx, filtered, updated)
	if err != nil {
		return err

	}

	return nil
}
