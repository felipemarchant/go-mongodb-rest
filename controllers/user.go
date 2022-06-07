package controllers

import (
	"github.com/felipemarchant/go-mongo-rest/database"
	"github.com/felipemarchant/go-mongo-rest/models"
	r "github.com/felipemarchant/go-mongo-rest/rest"
	"github.com/felipemarchant/go-mongo-rest/security"
	"github.com/felipemarchant/go-mongo-rest/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		r.Response(c, err.Error(), http.StatusBadRequest)
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
		r.Response(c, err.Error(), http.StatusInternalServerError)
		return
	}

	if count > 0 {
		r.Response(c, "User já existe.", http.StatusBadRequest)
		return
	}

	count, err = users.CountDocuments(ctx, bson.M{"phone": user.Phone})
	if err != nil {
		r.Response(c, err.Error(), http.StatusInternalServerError)
		return
	}

	if count > 0 {
		r.Response(c, "Phone está em uso.", http.StatusBadRequest)
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
		r.Response(c, "Houve um problema na inscrição do usuário.", http.StatusInternalServerError)
		return
	}

	r.Response(c, user, http.StatusCreated)
}
