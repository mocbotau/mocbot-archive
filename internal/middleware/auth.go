package middleware

import (
	"crypto/sha512"
	"encoding/hex"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/mocbotau/mocbot-archive/internal/utils"
)

// EnsureToken is a middleware that ensures the request has a valid API token.
func EnsureToken(secretManager *utils.SecretManager, secretPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.Request.Header.Get("X-API-Key")

		secret, err := secretManager.GetSecretFromFile(secretPath)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Internal server error"})
			log.Printf("Failed to read API secret: %v", err)

			return
		}

		hashBytes := sha512.Sum512([]byte(apiKey))

		hashString := hex.EncodeToString(hashBytes[:])
		if secret != hashString {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		}
	}
}
