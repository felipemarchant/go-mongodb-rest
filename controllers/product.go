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
	"net/http"
	"time"
)

func AddProduct(c *gin.Context) {
	var ctx, cancel = utils.ContextWithTimeout()
	var product models.Product
	defer cancel()

	if err := c.BindJSON(&product); err != nil {
		r.Response(c, err.Error(), http.StatusBadRequest)
		return
	}

	products := database.Client.ProductCollection()

	product.Id = primitive.NewObjectID()
	_, err := products.InsertOne(ctx, product)
	if err != nil {
		r.Response(c, "Não possível adicionar o produto", http.StatusInternalServerError)
		return
	}

	r.Response(c, "Produto adicionado com sucesso", http.StatusOK)
}

func GetProducts(c *gin.Context) {
	productList := make([]models.Product, 0)
	var ctx, cancel = utils.ContextWithTimeout()
	defer cancel()

	products := database.Client.ProductCollection()

	cursor, err := products.Find(ctx, bson.D{{}})
	if err != nil {
		r.Response(c, "Houve um problema ao buscar o produto", http.StatusInternalServerError)
		return
	}

	err = cursor.All(ctx, &productList)
	if err != nil {
		r.Response(c, err.Error(), http.StatusInternalServerError)
		return
	}

	defer cursor.Close(ctx)
	if err := cursor.Err(); err != nil {
		r.Response(c, err.Error(), http.StatusBadRequest)
		return
	}

	r.Response(c, productList, http.StatusOK)
}

func SearchProduct(c *gin.Context) {
	searchProducts := make([]models.Product, 0)
	queryParam := c.Query("name")
	if queryParam == "" {
		r.Response(c, "Critério de busca inválido", http.StatusBadRequest)
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	products := database.Client.ProductCollection()

	searchQuery, err := products.Find(ctx, bson.M{"name": bson.M{"$regex": queryParam}})
	if err != nil {
		r.Response(c, "Houve um problema ao buscar os produtos", http.StatusInternalServerError)
		return
	}

	err = searchQuery.All(ctx, &searchProducts)
	if err != nil {
		r.Response(c, err.Error(), http.StatusBadRequest)
		return
	}

	defer searchQuery.Close(ctx)
	if err := searchQuery.Err(); err != nil {
		r.Response(c, err.Error(), http.StatusBadRequest)
		return
	}

	r.Response(c, searchProducts, http.StatusOK)
}
