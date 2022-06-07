package controllers

import (
	"github.com/felipemarchant/go-mongo-rest/database"
	"github.com/felipemarchant/go-mongo-rest/models"
	"github.com/felipemarchant/go-mongo-rest/security"
	"github.com/felipemarchant/go-mongo-rest/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"time"
)

var Validate = validator.New()

func SignUp(c *gin.Context) {
	var ctx, cancel = utils.ContextWithTimeout()
	defer cancel()
	defer ctx.Done()

	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationErr := Validate.Struct(user)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	users := database.Client.UserCollection()

	count, err := users.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "data": err.Error()})
		log.Panic(err)
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "data": "User já existe."})
		return
	}

	count, err = users.CountDocuments(ctx, bson.M{"phone": user.Phone})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "data": err.Error()})
		log.Panic(err)
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": http.StatusBadRequest, "data": "Phone está em uso."})
		return
	}

	password := security.HashPassword(*user.Password)
	user.Password = &password
	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.ID = primitive.NewObjectID()
	user.UserId = user.ID.Hex()
	token, refreshToken, _ := security.TokenGenerator(*user.Email, *user.FirstName, *user.LastName, user.UserId)
	user.Token = &token
	user.RefreshToken = &refreshToken
	user.UserCart = make([]models.ProductUser, 0)
	user.AddressDetails = make([]models.Address, 0)
	user.OrderStatus = make([]models.Order, 0)

	_, insertErr := users.InsertOne(ctx, user)
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": http.StatusInternalServerError, "data": "Houve um problema na inscrição do usuário."})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "data": user})
}
