package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"task_manager/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func jwtSecret() string {
	return strings.TrimSpace(os.Getenv("JWT_SECRET"))
}

func GenerateToken(user models.User, ttl time.Duration) (string, error) {
	secret := jwtSecret()
	if secret == "" {
		return "", errors.New("JWT_SECRET is required")
	}

	claims := AuthClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     models.NormalizeRole(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		secret := jwtSecret()
		if secret == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "jwt secret is not configured"})
			c.Abort()
			return
		}

		claims := &AuthClaims{}
		token, err := jwt.ParseWithClaims(parts[1], claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("authClaims", claims)
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", models.NormalizeRole(claims.Role))

		c.Next()
	}
}

func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")

		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "role not found"})
			c.Abort()
			return
		}

		currentRole, ok := role.(string)
		if !ok || models.NormalizeRole(currentRole) != models.NormalizeRole(requiredRole) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}
		c.Next()
	}
}
