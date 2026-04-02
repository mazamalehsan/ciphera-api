package middlewares

import (
	"ciphera-api/env"
	"ciphera-api/jwt"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		pubKeyPEM := env.GetJWTPublicKey()
		pubKeyPEM = strings.ReplaceAll(pubKeyPEM, `\n`, "\n")
		block, _ := pem.Decode([]byte(pubKeyPEM))
		if block == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid server key config"})
			return
		}

		pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid server key config"})
			return
		}

		rsaPubKey, ok := pubKey.(*rsa.PublicKey)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid server key config"})
			return
		}

		claims, err := jwt.VerifyAccessToken(tokenStr, rsaPubKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("userId", claims.UserID)
		c.Next()
	}
}
