package controllers

import (
	"context"
	"fmt"
	"goproject/database"
	"goproject/helpers"
	"goproject/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	check := true
	msg := ""

	if err != nil {
		log.Panic(err)
		msg = fmt.Sprintf("Login or password incorrect.")
		check = false
	}

	return check, msg
}

func SignUp() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		if err := ginContext.BindJSON(&user); err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationError := validate.Struct(user)
		if validationError != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}

		if count > 0 {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exists"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, user.User_id)
		user.Token = &token
		user.RefreshToken = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		ginContext.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User

		if err := ginContext.BindJSON(&user); err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Login or passowrd is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*foundUser.Password, *user.Password)
		defer cancel()
		if passwordIsValid != true {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *&foundUser.User_id)

		helpers.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		ginContext.JSON(http.StatusOK, foundUser)
	}
}
