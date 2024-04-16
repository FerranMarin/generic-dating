package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/FerranMarin/generic-dating/initializers"
	"github.com/FerranMarin/generic-dating/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	authHeader := c.Request.Header["Authorization"]
	if authHeader == nil || !strings.HasPrefix(authHeader[0], "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Could not Authenticate"})
		return
	}

	// Extract the tokenString part without the "Bearer " prefix
	tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")

	// Decode/validate jwt
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Could not Authenticate"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Check token is not expired
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Could not Authenticate. Token is Expired"})
			return
		}
		// Find user and attach it to request
		var user models.User
		initializers.DB.First(&user, claims["sub"])
		if user.ID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Could not Authenticate"})
			return
		}
		c.Set("user", user)
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Could not Authenticate"})
		return
	}
}
