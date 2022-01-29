package main

import (
	"goproject/middleware"
	"goproject/routes"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

type User struct {
	ID uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var user = User {
	ID: 1,
	Username: "username",
	Password: "password",
}

func Login(context *gin.Context) {
	var requestUser User
	if err := context.ShouldBindJSON(&requestUser); err != nil {
		context.JSON(http.StatusUnprocessableEntity, "Invalid JSON provided")
		return
	}

	if requestUser.Username != user.Username || requestUser.Password != user.Password {
		context.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}

	token, err := CreateToken(user.ID)
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	context.JSON(http.StatusOK, token)
}

func CreateToken(userid uint64) (string, error) {
	var err error
	//Creating access token
	os.Setenv("ACCESS_SECRET", "jzdqnzqdnq")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	authenticationToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := authenticationToken.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}

	return token, nil
}




func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.UserRoutes(router)

	router.Use(middleware.Authentication())

	// API - 1
	router.GET("/api-1", func(ginContext *gin.Context) {
		ginContext.JSON(200, gin.H{"success": "Access granted for api-1"})
	})

	// API - 2
	router.GET("/api-2", func(ginContext *gin.Context) {
		ginContext.JSON(200, gin.H{"success": "Access granted for api-2"})
	})

	router.Run(":" + port)
}
