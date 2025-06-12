package main

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret = []byte("secret_key")
	tokenUses = make(map[string]int)
	mutex     sync.Mutex
)

type TokenClaims struct {
	jwt.RegisteredClaims
	TokenID string `json:"token_id"`
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	// Routes
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.POST("/generate-token", generateTokenHandler)
	r.GET("/characters", authMiddleware(), getCharactersHandler)

	r.Run(":8080")
}
