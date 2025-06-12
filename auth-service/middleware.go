package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Autorizacion en header es requerida"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token invalido o expirado"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*TokenClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		mutex.Lock()
		uses, exists := tokenUses[claims.TokenID]
		if !exists {
			mutex.Unlock()
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token no encontrado"})
			c.Abort()
			return
		}

		if uses >= 5 {
			mutex.Unlock()
			c.JSON(http.StatusForbidden, gin.H{"error": "Token ha expirado (uso maximo)"})
			c.Abort()
			return
		}

		tokenUses[claims.TokenID] = uses + 1
		mutex.Unlock()

		c.Set("tokenID", claims.TokenID)
		c.Next()
	}
}
