// internal/middleware/auth.go
package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Dubjay18/ecom-api/pkg/common/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("AuthMiddleware")
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header required", nil)
			return
		}

		// Check if the header has the Bearer prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "invalid authorization header format", nil)
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				response.Error(c, http.StatusUnauthorized, "token has expired", nil)
				return
			}
			response.Error(c, http.StatusUnauthorized, "invalid  token", nil)
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Add user ID and role to context
			if userID, ok := claims["user_id"]; ok && userID != nil {
				c.Set("userID", uint(userID.(float64)))
			} else {
				response.Error(c, http.StatusUnauthorized, "invalid  token", nil)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id in token"})
				return
			}
			if isAdmin, ok := claims["is_admin"]; ok && isAdmin != nil {
				c.Set("isAdmin", isAdmin.(bool))
			} else {
				response.Error(c, http.StatusUnauthorized, "invalid  token", nil)
				return
			}
			if email, ok := claims["email"]; ok && email != nil {
				c.Set("email", email.(string))
			} else {
				response.Error(c, http.StatusUnauthorized, "invalid token", nil)
				return
			}
			c.Next()
		} else {
			response.Error(c, http.StatusUnauthorized, "invalid token", nil)
			return
		}
	}
}
