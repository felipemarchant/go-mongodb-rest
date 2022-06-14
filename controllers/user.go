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
		r.Response(c, validationErr.Error(), http.StatusBadRequest)
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
	user.Id = primitive.NewObjectID()
	token, refreshToken, _ := security.TokenGenerator(*user.Email, *user.FirstName, *user.LastName, user.Id)
	user.Token = &token
	user.RefreshToken = &refreshToken
	user.UserCart = make([]models.Cart, 0)
	user.Addresses = make([]models.Address, 0)
	user.Orders = make([]models.Order, 0)
	_, insertErr := users.InsertOne(ctx, user)
	if insertErr != nil {
		r.Response(c, "Houve um problema na inscrição do usuário.", http.StatusInternalServerError)
		return
	}
	r.Response(c, user, http.StatusCreated)
}

func Login(c *gin.Context) {
	var ctx, cancel = utils.ContextWithTimeout()
	defer cancel()
	var user models.User
	var foundUser models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	users := database.Client.UserCollection()
	err := users.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
	if err != nil {
		r.Response(c, "Login ou Password incorreto", http.StatusUnauthorized)
		return
	}
	PasswordIsValid, msg := security.VerifyPassword(*user.Password, *foundUser.Password)
	if !PasswordIsValid {
		r.Response(c, msg, http.StatusInternalServerError)
		return
	}
	token, refreshToken, _ := security.TokenGenerator(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, foundUser.Id)
	security.UpdateAllTokens(token, refreshToken, foundUser.Id.String())
	r.Response(c, foundUser, http.StatusFound)
}
