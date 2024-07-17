package middleware

import (
	"my_module/api/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Check(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")

	if accessToken == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Authorization is required",
		})
		return
	}

	_, err := auth.ValidateAccessToken(accessToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	c.Next()
}
