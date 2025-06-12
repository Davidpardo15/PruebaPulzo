package main

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func generateTokenHandler(c *gin.Context) {
	tokenID := generateRandomID()
	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		TokenID: tokenID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al generar el token"})
		return
	}

	mutex.Lock()
	tokenUses[tokenID] = 0
	mutex.Unlock()

	c.JSON(http.StatusOK, gin.H{"token": signedToken})
}

func getCharactersHandler(c *gin.Context) {
	// This will be called after the authMiddleware validates the token

	resp, err := http.Get("http://localhost:8081/api/characters")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al traer los personajes"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "Data service error"})
		return
	}

	// Read response body and forward it
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	c.Data(resp.StatusCode, "application/json", body)
}

func generateRandomID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
