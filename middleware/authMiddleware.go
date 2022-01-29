package middleware

import (
    "fmt"
    "net/http"

    helper "goproject/helpers"

    "github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		clientToken := ginContext.Request.Header.Get("token")
		if clientToken == "" {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			ginContext.Abort()
			return
		}

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err})
			ginContext.Abort()
			return
		}

		ginContext.Set("email", claims.Email)
		ginContext.Set("uid", claims.Uid)

		ginContext.Next()
	}
}