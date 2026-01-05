package middleware

import (
	"strings"
	"users/auth"
	"users/users"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AuthMiddleware(jwtService auth.Service, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "missing token"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 1️⃣ validate JWT
		token, err := jwtService.ValidateToken(
			tokenString,
			c.Request.Context(),
		)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid token"})
			return
		}

		// 2️⃣ ambil claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid claims"})
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok || userIDFloat == 0 {
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid user_id"})
			return
		}

		userID := uint(userIDFloat)

		// 3️⃣ ambil user dari DB
		var user users.Users
		if err := db.First(&user, userID).Error; err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "user not found"})
			return
		}

		// 4️⃣ set ke context
		c.Set("user", user)
		c.Next()
	}
}
