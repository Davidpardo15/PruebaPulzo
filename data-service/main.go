package main

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/api/characters", getCharacters)

	r.Run(":8081")
}

func getCharacters(c *gin.Context) {
	// Call Rick and Morty API
	resp, err := http.Get("https://rickandmortyapi.com/api/character")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch characters from API"})
		return
	}
	defer resp.Body.Close()

	// Forward the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read API response"})
		return
	}

	c.Data(resp.StatusCode, "application/json", body)
}
